package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"strings"

	"golang-ecommerce-search/internal/domain"

	"github.com/elastic/go-elasticsearch/v8"
)

type productRepository struct {
	client *elasticsearch.Client
	index  string
}

type ProductRepository interface {
	Search(params domain.SearchParams) ([]*domain.Product, error)
	Create(product *domain.Product) error
	Update(product *domain.Product) error
	Delete(id string) error
	GetByID(id string) (*domain.Product, error)
	IncrementViews(id string) error
	IncrementBuys(id string) error
}

func NewProductRepository(client *elasticsearch.Client, index string) ProductRepository {
	return &productRepository{
		client: client,
		index:  index,
	}
}

func (r *productRepository) Create(product *domain.Product) error {
	ctx := context.Background()
	body, err := json.Marshal(product)
	if err != nil {
		return err
	}

	_, err = r.client.Index(
		r.index,
		bytes.NewReader(body),
		r.client.Index.WithContext(ctx),
		r.client.Index.WithDocumentID(product.ID),
	)
	return err
}

func (r *productRepository) Update(product *domain.Product) error {
	ctx := context.Background()
	body, err := json.Marshal(product)
	if err != nil {
		return err
	}

	_, err = r.client.Update(
		r.index,
		product.ID,
		bytes.NewReader([]byte(`{"doc":`+string(body)+`}`)),
		r.client.Update.WithContext(ctx),
	)
	return err
}

func (r *productRepository) Delete(id string) error {
	ctx := context.Background()
	_, err := r.client.Delete(
		r.index,
		id,
		r.client.Delete.WithContext(ctx),
	)
	return err
}

func (r *productRepository) GetByID(id string) (*domain.Product, error) {
	ctx := context.Background()
	res, err := r.client.Get(
		r.index,
		id,
		r.client.Get.WithContext(ctx),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result struct {
		Source domain.Product `json:"_source"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Source, nil
}

func (r *productRepository) IncrementViews(id string) error {
	ctx := context.Background()
	update := map[string]interface{}{
		"script": map[string]interface{}{
			"source": "ctx._source.views += 1",
			"lang":   "painless",
		},
	}

	body, err := json.Marshal(update)
	if err != nil {
		return err
	}

	_, err = r.client.Update(
		r.index,
		id,
		bytes.NewReader(body),
		r.client.Update.WithContext(ctx),
	)
	return err
}

func (r *productRepository) IncrementBuys(id string) error {
	ctx := context.Background()
	update := map[string]interface{}{
		"script": map[string]interface{}{
			"source": "ctx._source.buys += 1",
			"lang":   "painless",
		},
	}

	body, err := json.Marshal(update)
	if err != nil {
		return err
	}

	_, err = r.client.Update(
		r.index,
		id,
		bytes.NewReader(body),
		r.client.Update.WithContext(ctx),
	)
	return err
}

func (r *productRepository) Search(params domain.SearchParams) ([]*domain.Product, error) {
	ctx := context.Background()
	query := strings.ToLower(params.Query)

	// Build the query
	queryMap := map[string]interface{}{
		"bool": map[string]interface{}{
			"must": []map[string]interface{}{},
		},
	}

	// Add text search if query is provided
	if query != "" {
		queryMap["bool"].(map[string]interface{})["must"] = append(
			queryMap["bool"].(map[string]interface{})["must"].([]map[string]interface{}),
			map[string]interface{}{
				"multi_match": map[string]interface{}{
					"query":  query,
					"fields": []string{"name^3", "description^2", "category", "tags"},
				},
			},
		)
	}

	// Add category filter if provided
	if len(params.Categories) > 0 {
		queryMap["bool"].(map[string]interface{})["must"] = append(
			queryMap["bool"].(map[string]interface{})["must"].([]map[string]interface{}),
			map[string]interface{}{
				"terms": map[string]interface{}{
					"category.keyword": params.Categories,
				},
			},
		)
	}

	// Add brand filter if provided
	if len(params.Brands) > 0 {
		queryMap["bool"].(map[string]interface{})["must"] = append(
			queryMap["bool"].(map[string]interface{})["must"].([]map[string]interface{}),
			map[string]interface{}{
				"terms": map[string]interface{}{
					"brand.keyword": params.Brands,
				},
			},
		)
	}

	// Build the sort
	var sort []map[string]interface{}
	switch params.SortBy {
	case "views":
		sort = append(sort, map[string]interface{}{"views": "desc"})
	case "buys":
		sort = append(sort, map[string]interface{}{"buys": "desc"})
	default:
		// Default sort by score (relevance) and then by views and buys
		sort = append(sort, map[string]interface{}{"_score": "desc"})
		sort = append(sort, map[string]interface{}{"views": "desc"})
		sort = append(sort, map[string]interface{}{"buys": "desc"})
	}

	// Calculate pagination
	from := (params.Page - 1) * params.PageSize
	if from < 0 {
		from = 0
	}

	body := map[string]interface{}{
		"query": map[string]interface{}{
			"function_score": map[string]interface{}{
				"query": queryMap,
				"functions": []map[string]interface{}{
					{
						"field_value_factor": map[string]interface{}{
							"field":    "buys",
							"factor":   0.3,
							"modifier": "log1p",
						},
					},
					{
						"field_value_factor": map[string]interface{}{
							"field":    "views",
							"factor":   0.1,
							"modifier": "log1p",
						},
					},
				},
				"score_mode": "sum",
				"boost_mode": "sum",
			},
		},
		"sort": sort,
		"from": from,
		"size": params.PageSize,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Log the query for debugging
	log.Printf("Elasticsearch Query: %s\n", string(bodyBytes))

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(r.index),
		r.client.Search.WithBody(bytes.NewReader(bodyBytes)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result struct {
		Hits struct {
			Hits []struct {
				Source domain.Product `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	products := make([]*domain.Product, len(result.Hits.Hits))
	for i, hit := range result.Hits.Hits {
		products[i] = &hit.Source
	}

	return products, nil
}
