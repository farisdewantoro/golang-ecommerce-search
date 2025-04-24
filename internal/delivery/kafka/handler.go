package kafka

import (
	"encoding/json"

	"golang-ecommerce-search/internal/domain"
)

type ProductEventHandler struct {
	productService domain.ProductService
}

func NewProductEventHandler(productService domain.ProductService) *ProductEventHandler {
	return &ProductEventHandler{
		productService: productService,
	}
}

func (h *ProductEventHandler) OnCreated(message []byte) error {
	var product domain.Product
	if err := json.Unmarshal(message, &product); err != nil {
		return err
	}

	return h.productService.OnCreated(&product)
}

func (h *ProductEventHandler) OnUpdated(message []byte) error {
	var product domain.Product
	if err := json.Unmarshal(message, &product); err != nil {
		return err
	}

	return h.productService.OnUpdated(&product)
}

func (h *ProductEventHandler) OnDeleted(message []byte) error {
	productID := string(message)
	return h.productService.OnDeleted(productID)
}

func (h *ProductEventHandler) OnViewsIncremented(message []byte) error {
	productID := string(message)
	return h.productService.OnViewsIncremented(productID)
}

func (h *ProductEventHandler) OnBuysIncremented(message []byte) error {
	productID := string(message)
	return h.productService.OnBuysIncremented(productID)
}
