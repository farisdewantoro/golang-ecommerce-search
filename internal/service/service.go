package service

import (
	"golang-ecommerce-search/internal/config"
	"golang-ecommerce-search/internal/domain"
	es "golang-ecommerce-search/internal/repository/elasticsearch"
	mongo "golang-ecommerce-search/internal/repository/mongodb"
	"golang-ecommerce-search/pkg/kafka"
)

type ProductService interface {
	domain.ProductService
}

type productService struct {
	esRepo    es.ProductRepository
	mongoRepo mongo.ProductRepository
	producer  *kafka.Producer
	config    *config.Config
}

func NewProductService(esRepo es.ProductRepository, mongoRepo mongo.ProductRepository, producer *kafka.Producer, cfg *config.Config) ProductService {
	return &productService{
		esRepo:    esRepo,
		mongoRepo: mongoRepo,
		producer:  producer,
		config:    cfg,
	}
}
