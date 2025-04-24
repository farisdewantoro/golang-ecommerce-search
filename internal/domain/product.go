package domain

import "time"

type Product struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`
	Price       float64   `json:"price" bson:"price"`
	Category    string    `json:"category" bson:"category"`
	Tags        []string  `json:"tags" bson:"tags"`
	Brand       string    `json:"brand" bson:"brand"`
	Views       int64     `json:"views" bson:"views"`
	Buys        int64     `json:"buys" bson:"buys"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

type SearchParams struct {
	Query      string
	Categories []string
	Brands     []string
	SortBy     string
	Page       int
	PageSize   int
}

type ProductService interface {
	CreateProduct(product *Product) error
	UpdateProduct(product *Product) error
	DeleteProduct(id string) error
	GetProduct(id string) (*Product, error)
	SearchProducts(params SearchParams) ([]*Product, error)
	IncrementViews(id string) error
	IncrementBuys(id string) error
	OnCreated(product *Product) error
	OnUpdated(product *Product) error
	OnDeleted(productID string) error
	OnViewsIncremented(productID string) error
	OnBuysIncremented(productID string) error
}
