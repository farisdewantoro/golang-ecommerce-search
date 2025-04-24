package main

import (
	"log"

	"golang-ecommerce-search/internal/config"
	"golang-ecommerce-search/internal/delivery/http/handler"
	"golang-ecommerce-search/internal/repository/elasticsearch"
	"golang-ecommerce-search/internal/repository/mongodb"
	"golang-ecommerce-search/internal/service"
	"golang-ecommerce-search/pkg/esclient"
	"golang-ecommerce-search/pkg/kafka"
	"golang-ecommerce-search/pkg/mongodbclient"

	"github.com/gin-gonic/gin"
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

	// Initialize Kafka producer
	kafkaProducer, err := kafka.NewProducer(&kafka.Config{
		Brokers: cfg.Kafka.Brokers,
	})
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer kafkaProducer.Close()

	// Initialize repositories and services
	productRepo := mongodb.NewProductRepository(mongoClient.GetDatabase(), cfg.MongoDB.Collection)
	esRepo := elasticsearch.NewProductRepository(esClient.GetClient(), cfg.Elasticsearch.Index)
	productService := service.NewProductService(esRepo, productRepo, kafkaProducer, cfg)
	productHandler := handler.NewProductHandler(productService)

	// Initialize Gin router
	router := gin.Default()

	// Register routes
	router.POST("/products", productHandler.Create)
	router.PUT("/products/:id", productHandler.Update)
	router.DELETE("/products/:id", productHandler.Delete)
	router.GET("/products/:id", productHandler.Get)
	router.GET("/products/search", productHandler.Search)
	router.POST("/products/:id/views", productHandler.IncrementViews)
	router.POST("/products/:id/buys", productHandler.IncrementBuys)

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
