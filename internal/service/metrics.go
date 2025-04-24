package service

import (
	"fmt"
)

func (s *productService) IncrementViews(id string) error {
	if err := s.mongoRepo.IncrementViews(id); err != nil {
		return fmt.Errorf("failed to increment views in MongoDB: %w", err)
	}

	return s.publishEvent(s.config.Kafka.Topic.ProductViewsInc, id)
}

func (s *productService) IncrementBuys(id string) error {
	if err := s.mongoRepo.IncrementBuys(id); err != nil {
		return fmt.Errorf("failed to increment buys in MongoDB: %w", err)
	}

	return s.publishEvent(s.config.Kafka.Topic.ProductBuysInc, id)
}
