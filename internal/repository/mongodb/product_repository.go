package mongodb

import (
	"context"
	"fmt"
	"time"

	"golang-ecommerce-search/internal/domain"
	"golang-ecommerce-search/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepository interface {
	Create(product *domain.Product) error
	Update(product *domain.Product) error
	Delete(id string) error
	GetByID(id string) (*domain.Product, error)
	Search(params domain.SearchParams) ([]*domain.Product, error)
	IncrementViews(id string) error
	IncrementBuys(id string) error
}

type productRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database, collectionName string) ProductRepository {
	collection := db.Collection(collectionName)
	return &productRepository{
		collection: collection,
	}
}

func (r *productRepository) Create(product *domain.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Generate new UUID for the product
	productID := model.NewID()
	product.ID = productID.String()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, product)
	return err
}

func (r *productRepository) Update(product *domain.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	product.UpdatedAt = time.Now()

	filter := bson.M{"_id": product.ID}
	update := bson.M{
		"$set": bson.M{
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"category":    product.Category,
			"brand":       product.Brand,
			"tags":        product.Tags,
			"updated_at":  product.UpdatedAt,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *productRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

func (r *productRepository) GetByID(id string) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var product domain.Product
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Search(params domain.SearchParams) ([]*domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build the filter
	filter := bson.M{}

	// Add text search if query is provided
	if params.Query != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": params.Query, "$options": "i"}},
			{"description": bson.M{"$regex": params.Query, "$options": "i"}},
			{"category": bson.M{"$regex": params.Query, "$options": "i"}},
			{"tags": bson.M{"$in": []string{params.Query}}},
		}
	}

	// Add category filter if provided
	if len(params.Categories) > 0 {
		filter["category"] = bson.M{"$in": params.Categories}
	}

	// Add brand filter if provided
	if len(params.Brands) > 0 {
		filter["brand"] = bson.M{"$in": params.Brands}
	}

	// Build the sort options
	sort := bson.M{}
	switch params.SortBy {
	case "views":
		sort["views"] = -1
	case "buys":
		sort["buys"] = -1
	default:
		// Default sort by views and buys
		sort["views"] = -1
		sort["buys"] = -1
	}

	// Calculate pagination
	skip := (params.Page - 1) * params.PageSize
	if skip < 0 {
		skip = 0
	}

	// Find products with pagination and sorting
	cursor, err := r.collection.Find(ctx, filter, options.Find().
		SetSort(sort).
		SetSkip(int64(skip)).
		SetLimit(int64(params.PageSize)))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*domain.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) IncrementViews(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	update := bson.M{
		"$inc": bson.M{
			"views": 1,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	// Check if document was actually updated
	if result.MatchedCount == 0 {
		return fmt.Errorf("product with ID %s not found", id)
	}

	return nil
}

func (r *productRepository) IncrementBuys(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	update := bson.M{
		"$inc": bson.M{
			"buys": 1,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	// Check if document was actually updated
	if result.MatchedCount == 0 {
		return fmt.Errorf("product with ID %s not found", id)
	}

	return nil
}
