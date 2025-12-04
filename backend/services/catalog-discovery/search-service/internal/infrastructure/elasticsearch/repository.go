package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/titan-commerce/backend/search-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type SearchRepository struct {
	client *elasticsearch.Client
	index  string
	logger *logger.Logger
}

func NewSearchRepository(addresses []string, logger *logger.Logger) (*SearchRepository, error) {
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to create elasticsearch client", err)
	}

	// Ping to verify connection
	res, err := client.Info()
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to connect to elasticsearch", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New(errors.ErrInternal, "elasticsearch returned error")
	}

	logger.Info("Elasticsearch repository initialized")
	return &SearchRepository{
		client: client,
		index:  "products",
		logger: logger,
	}, nil
}

func (r *SearchRepository) IndexProduct(ctx context.Context, product *domain.ProductDocument) error {
	data, err := json.Marshal(product)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to marshal product", err)
	}

	req := esapi.IndexRequest{
		Index:      r.index,
		DocumentID: product.ID,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to index product", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return errors.New(errors.ErrInternal, "elasticsearch indexing failed")
	}

	return nil
}

func (r *SearchRepository) Search(ctx context.Context, query string, page, pageSize int) ([]*domain.ProductDocument, int, error) {
	from := (page - 1) * pageSize
	
	// Simple multi-match query
	searchBody := map[string]interface{}{
		"from": from,
		"size": pageSize,
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"name^3", "description", "category"},
				"fuzziness": "AUTO",
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(searchBody); err != nil {
		return nil, 0, errors.Wrap(errors.ErrInternal, "failed to encode search query", err)
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(r.index),
		r.client.Search.WithBody(&buf),
		r.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, 0, errors.Wrap(errors.ErrInternal, "failed to execute search", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, 0, errors.New(errors.ErrInternal, "elasticsearch search failed")
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, 0, errors.Wrap(errors.ErrInternal, "failed to decode search response", err)
	}

	hits := result["hits"].(map[string]interface{})
	total := int(hits["total"].(map[string]interface{})["value"].(float64))
	hitsList := hits["hits"].([]interface{})

	var products []*domain.ProductDocument
	for _, hit := range hitsList {
		source := hit.(map[string]interface{})["_source"]
		sourceBytes, _ := json.Marshal(source)
		var p domain.ProductDocument
		json.Unmarshal(sourceBytes, &p)
		products = append(products, &p)
	}

	return products, total, nil
}
