package service

import (
	"fmt"

	"golang-ecommerce-search/internal/domain"
)

func (s *productService) CreateProduct(product *domain.Product) error {
	if err := s.mongoRepo.Create(product); err != nil {
		return fmt.Errorf("failed to create product in MongoDB: %w", err)
	}

	return s.publishEvent(s.config.Kafka.Topic.ProductCreated, product)
}

func (s *productService) UpdateProduct(product *domain.Product) error {
	if err := s.mongoRepo.Update(product); err != nil {
		return fmt.Errorf("failed to update product in MongoDB: %w", err)
	}

	return s.publishEvent(s.config.Kafka.Topic.ProductUpdated, product)
}

func (s *productService) DeleteProduct(id string) error {
	if err := s.mongoRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete product from MongoDB: %w", err)
	}

	return s.publishEvent(s.config.Kafka.Topic.ProductDeleted, id)
}

func (s *productService) GetProduct(id string) (*domain.Product, error) {
	product, err := s.mongoRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product from MongoDB: %w", err)
	}
	return product, nil
}

func (s *productService) SearchProducts(params domain.SearchParams) ([]*domain.Product, error) {
	products, err := s.esRepo.Search(params)
	if err != nil {
		return nil, fmt.Errorf("failed to search products in Elasticsearch: %w", err)
	}
	return products, nil
}
