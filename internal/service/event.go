package service

import (
	"encoding/json"
	"fmt"

	"golang-ecommerce-search/internal/domain"
)

// publishEvent publishes an event to Kafka with the given topic and payload
func (s *productService) publishEvent(topic string, payload interface{}) error {
	var message string
	switch p := payload.(type) {
	case string:
		message = p
	default:
		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal event payload: %w", err)
		}
		message = string(jsonBytes)
	}

	if err := s.producer.SendMessage(topic, message); err != nil {
		return fmt.Errorf("failed to publish event to topic %s: %w", topic, err)
	}
	return nil
}

// Event handlers for Elasticsearch synchronization
func (s *productService) OnCreated(product *domain.Product) error {
	if err := s.esRepo.Create(product); err != nil {
		return fmt.Errorf("failed to create product in Elasticsearch: %w", err)
	}
	return nil
}

func (s *productService) OnUpdated(product *domain.Product) error {
	if err := s.esRepo.Update(product); err != nil {
		return fmt.Errorf("failed to update product in Elasticsearch: %w", err)
	}
	return nil
}

func (s *productService) OnDeleted(productID string) error {
	if err := s.esRepo.Delete(productID); err != nil {
		return fmt.Errorf("failed to delete product from Elasticsearch: %w", err)
	}
	return nil
}

func (s *productService) OnViewsIncremented(productID string) error {
	product, err := s.mongoRepo.GetByID(productID)
	if err != nil {
		return fmt.Errorf("failed to get product for views increment sync: %w", err)
	}
	return s.OnUpdated(product)
}

func (s *productService) OnBuysIncremented(productID string) error {
	product, err := s.mongoRepo.GetByID(productID)
	if err != nil {
		return fmt.Errorf("failed to get product for buys increment sync: %w", err)
	}
	return s.OnUpdated(product)
}
