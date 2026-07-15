package elasticsearch

import (
	"context"
	"fmt"

	"github.com/artuos/sniffer/internal/domain"
	"github.com/olivere/elastic/v7"
)

const indexName = "fragrances"

var fragranceMapping = `{
  "mappings": {
    "properties": {
      "name":         { "type": "text", "fields": { "keyword": { "type": "keyword" } } },
      "gender":       { "type": "keyword" },
      "rating_value": { "type": "float" },
      "rating_count": { "type": "integer" },
      "main_accords": { "type": "keyword" },
      "perfumers":    { "type": "text", "fields": { "keyword": { "type": "keyword" } } },
      "description":  { "type": "text" },
      "url":          { "type": "keyword" }
    }
  }
}`

type ESFragranceRepositoryAdapter struct {
	db *ESClientAdapter
}

func NewESFragranceRepositoryAdapter(client *ESClientAdapter) *ESFragranceRepositoryAdapter {
	return &ESFragranceRepositoryAdapter{
		db: client,
	}
}

func (a *ESFragranceRepositoryAdapter) Create(ctx context.Context, fragrances []domain.Fragrance) error {
	bulk := a.db.Client.Bulk().Index(indexName)

	for _, doc := range fragrances {
		bulk.Add(elastic.NewBulkIndexRequest().Id(doc.ID.String()).Doc(doc))
		if bulk.NumberOfActions() >= 500 {
			if _, err := bulk.Do(ctx); err != nil {
				return fmt.Errorf("bulk index: %w", err)
			}
		}
	}

	if bulk.NumberOfActions() > 0 {
		if _, err := bulk.Do(ctx); err != nil {
			return fmt.Errorf("bulk index flush: %w", err)
		}
	}
	return nil
}
