package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
)

const indexName = "products"

type ProductDocument struct {
	ProductID   string  `json:"product_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	CategoryID  string  `json:"category_id"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
	Rating      float64 `json:"rating"`
	SoldCount   int     `json:"sold_count"`
	IndexedAt   string  `json:"indexed_at"`
}

type SearchResult struct {
	ProductID   string
	Name        string
	Description string
	Price       float64
	ImageURL    string
	Rating      float64
	SoldCount   int
	Score       float32
}

type ElasticsearchRepository struct {
	client *elasticsearch.Client
	logger *logger.Logger
}

func NewElasticsearchRepository(addresses []string, logger *logger.Logger) (*ElasticsearchRepository, error) {
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to create Elasticsearch client", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	res, err := client.Info(esapi.WithContext(ctx))
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to connect to Elasticsearch", err)
	}
	defer res.Body.Close()

	logger.Info("Connected to Elasticsearch")
	return &ElasticsearchRepository{client: client, logger: logger}, nil
}

// IndexProduct indexes a product document
func (r *ElasticsearchRepository) IndexProduct(ctx context.Context, doc *ProductDocument) error {
	doc.IndexedAt = time.Now().UTC().Format(time.RFC3339)
	
	data, err := json.Marshal(doc)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to marshal document", err)
	}

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: doc.ProductID,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to index document", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return errors.New(errors.ErrInternal, fmt.Sprintf("error indexing document: %s", res.Status()))
	}

	r.logger.Infof("Indexed product: %s", doc.ProductID)
	return nil
}

// Search performs full-text search with filters
func (r *ElasticsearchRepository) Search(ctx context.Context, query, categoryID string, minPrice, maxPrice, minRating float64, page, pageSize int) ([]SearchResult, int, error) {
	// Build Elasticsearch query
	var buf bytes.Buffer
	searchQuery := map[string]interface{}{
		"from": (page - 1) * pageSize,
		"size": pageSize,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"multi_match": map[string]interface{}{
							"query":  query,
							"fields": []string{"name^3", "description"},
						},
					},
				},
			},
		},
		"sort": []interface{}{
			map[string]interface{}{"_score": "desc"},
			map[string]interface{}{"sold_count": "desc"},
		},
	}

	// Add filters
	filters := []interface{}{}
	if categoryID != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{"category_id": categoryID},
		})
	}
	if minPrice > 0 || maxPrice > 0 {
		priceFilter := map[string]interface{}{}
		if minPrice > 0 {
			priceFilter["gte"] = minPrice
		}
		if maxPrice > 0 {
			priceFilter["lte"] = maxPrice
		}
		filters = append(filters, map[string]interface{}{
			"range": map[string]interface{}{"price": priceFilter},
		})
	}
	if minRating > 0 {
		filters = append(filters, map[string]interface{}{
			"range": map[string]interface{}{"rating": map[string]interface{}{"gte": minRating}},
		})
	}

	if len(filters) > 0 {
		searchQuery["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filters
	}

	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, 0, errors.Wrap(errors.ErrInternal, "failed to encode query", err)
	}

	// Execute search
	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(indexName),
		r.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, 0, errors.Wrap(errors.ErrInternal, "search request failed", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, 0, errors.New(errors.ErrInternal, fmt.Sprintf("search error: %s", res.Status()))
	}

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, 0, errors.Wrap(errors.ErrInternal, "failed to parse response", err)
	}

	hits := result["hits"].(map[string]interface{})
	totalHits := int(hits["total"].(map[string]interface{})["value"].(float64))
	
	results := []SearchResult{}
	for _, hit := range hits["hits"].([]interface{}) {
		h := hit.(map[string]interface{})
		source := h["_source"].(map[string]interface{})
		score := float32(h["_score"].(float64))

		results = append(results, SearchResult{
			ProductID:   source["product_id"].(string),
			Name:        source["name"].(string),
			Description: source["description"].(string),
			Price:       source["price"].(float64),
			ImageURL:    source["image_url"].(string),
			Rating:      source["rating"].(float64),
			SoldCount:   int(source["sold_count"].(float64)),
			Score:       score,
		})
	}

	return results, totalHits, nil
}

// DeleteProduct deletes a product from search index
func (r *ElasticsearchRepository) DeleteProduct(ctx context.Context, productID string) error {
	req := esapi.DeleteRequest{
		Index:      indexName,
		DocumentID: productID,
		Refresh:    "true",
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to delete document", err)
	}
	defer res.Body.Close()

	r.logger.Infof("Deleted product from search: %s", productID)
	return nil
}
