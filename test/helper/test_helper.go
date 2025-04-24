package helper

import (
	"context"
	"testing"
	"time"

	"golang-ecommerce-search/internal/config"
	"golang-ecommerce-search/pkg/esclient"
	"golang-ecommerce-search/pkg/kafka"
	"golang-ecommerce-search/pkg/mongodbclient"

	"github.com/stretchr/testify/require"
)

type TestEnv struct {
	Config    *config.Config
	MongoDB   *mongodbclient.Client
	ES        *esclient.Client
	KafkaProd *kafka.Producer
	KafkaCons *kafka.Consumer
}

func SetupTestEnv(t *testing.T) *TestEnv {
	// Load test configuration
	cfg, err := config.LoadConfig("../../config/config.test.yaml")
	require.NoError(t, err)

	// Initialize MongoDB
	mongoClient, err := mongodbclient.NewClient(&mongodbclient.Config{
		URI:      cfg.MongoDB.URI,
		Database: cfg.MongoDB.Database,
	})
	require.NoError(t, err)

	// Initialize Elasticsearch
	esClient, err := esclient.NewClient(&esclient.Config{
		Addresses: cfg.Elasticsearch.Addresses,
		Username:  cfg.Elasticsearch.Username,
		Password:  cfg.Elasticsearch.Password,
	})
	require.NoError(t, err)

	// Initialize Kafka
	kafkaProducer, err := kafka.NewProducer(&kafka.Config{
		Brokers: cfg.Kafka.Brokers,
	})
	require.NoError(t, err)

	kafkaConsumer, err := kafka.NewConsumer(&kafka.Config{
		Brokers: cfg.Kafka.Brokers,
	})
	require.NoError(t, err)

	return &TestEnv{
		Config:    cfg,
		MongoDB:   mongoClient,
		ES:        esClient,
		KafkaProd: kafkaProducer,
		KafkaCons: kafkaConsumer,
	}
}

func (env *TestEnv) Cleanup(t *testing.T) {
	// Clean up MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := env.MongoDB.GetDatabase().Drop(ctx)
	require.NoError(t, err)

	// Clean up Elasticsearch
	// Note: You might want to add specific cleanup for ES indices

	// Close connections
	env.MongoDB.Close()
	env.KafkaProd.Close()
	env.KafkaCons.Close()
}
