package main

import (
	"log"

	"golang-ecommerce-search/internal/config"
	"golang-ecommerce-search/internal/delivery/kafka"
	"golang-ecommerce-search/internal/repository/elasticsearch"
	"golang-ecommerce-search/internal/repository/mongodb"
	"golang-ecommerce-search/internal/service"
	"golang-ecommerce-search/pkg/esclient"
	kafkapkg "golang-ecommerce-search/pkg/kafka"
	"golang-ecommerce-search/pkg/mongodbclient"

	"github.com/Shopify/sarama"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize MongoDB client
	mongoClient, err := mongodbclient.NewClient(&mongodbclient.Config{
		URI:      cfg.MongoDB.URI,
		Database: cfg.MongoDB.Database,
	})
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Close()

	// Initialize Elasticsearch client
	esClient, err := esclient.NewClient(&esclient.Config{
		Addresses: cfg.Elasticsearch.Addresses,
		Username:  cfg.Elasticsearch.Username,
		Password:  cfg.Elasticsearch.Password,
	})
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %v", err)
	}

	// Initialize Kafka consumer
	kafkaConsumer, err := kafkapkg.NewConsumer(&kafkapkg.Config{
		Brokers: cfg.Kafka.Brokers,
		GroupID: cfg.Kafka.GroupID,
	})
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer kafkaConsumer.Close()

	// Initialize repositories
	mongoRepo := mongodb.NewProductRepository(mongoClient.GetDatabase(), cfg.MongoDB.Collection)
	esRepo := elasticsearch.NewProductRepository(esClient.GetClient(), cfg.Elasticsearch.Index)

	// Initialize product service
	productService := service.NewProductService(esRepo, mongoRepo, nil, cfg)

	// Initialize event handler
	eventHandler := kafka.NewProductEventHandler(productService)

	// Subscribe to Kafka topics
	topics := []string{
		cfg.Kafka.Topic.ProductCreated,
		cfg.Kafka.Topic.ProductUpdated,
		cfg.Kafka.Topic.ProductDeleted,
		cfg.Kafka.Topic.ProductViewsInc,
		cfg.Kafka.Topic.ProductBuysInc,
	}

	// Create consumers for each topic
	consumers := make([]sarama.PartitionConsumer, len(topics))
	for i, topic := range topics {
		consumer, err := kafkaConsumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
		if err != nil {
			log.Fatalf("Failed to create consumer for topic %s: %v", topic, err)
		}
		consumers[i] = consumer
		defer consumer.Close()
	}

	log.Println("Ready to consume messages...")

	// Process messages from all consumers
	for {
		for i, consumer := range consumers {
			select {
			case msg := <-consumer.Messages():
				log.Printf("Received message from topic %s: %s", topics[i], string(msg.Value))
				switch topics[i] {
				case cfg.Kafka.Topic.ProductCreated:
					if err := eventHandler.OnCreated(msg.Value); err != nil {
						log.Printf("Failed to handle product_created event: %v", err)
					}
				case cfg.Kafka.Topic.ProductUpdated:
					if err := eventHandler.OnUpdated(msg.Value); err != nil {
						log.Printf("Failed to handle product_updated event: %v", err)
					}
				case cfg.Kafka.Topic.ProductDeleted:
					if err := eventHandler.OnDeleted(msg.Value); err != nil {
						log.Printf("Failed to handle product_deleted event: %v", err)
					}
				case cfg.Kafka.Topic.ProductViewsInc:
					if err := eventHandler.OnViewsIncremented(msg.Value); err != nil {
						log.Printf("Failed to handle product_views_incremented event: %v", err)
					}
				case cfg.Kafka.Topic.ProductBuysInc:
					if err := eventHandler.OnBuysIncremented(msg.Value); err != nil {
						log.Printf("Failed to handle product_buys_incremented event: %v", err)
					}
				}
			default:
				// No message available, continue to next consumer
				continue
			}
		}
	}
}
