package elasticsearch

import (
	"context"
	"encoding/json"
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

func (a *ESFragranceRepositoryAdapter) Search(ctx context.Context, params domain.SearchParams) (*domain.SearchResponse, error) {
	q := elastic.NewBoolQuery()

	if params.Query != "" {
		q.Should(
			elastic.NewMatchPhrasePrefixQuery("name", params.Query).Boost(2.0),
			elastic.NewMultiMatchQuery(params.Query, "name^3", "perfumers", "main_accords").
				Type("best_fields").
				Fuzziness("AUTO"),
		).MinimumShouldMatch("1")
	}

	if params.Gender != "" {
		q.Filter(elastic.NewTermQuery("gender", params.Gender))
	}
	if params.Accord != "" {
		q.Filter(elastic.NewTermQuery("main_accords", params.Accord))
	}

	search := a.db.Client.Search().
		Index(indexName).
		Query(q).
		Size(20).
		Aggregation("by_gender", elastic.NewTermsAggregation().Field("gender")).
		Aggregation("by_main_accords", elastic.NewTermsAggregation().Field("main_accords"))

	result, err := search.Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}

	fragrances := make([]domain.Fragrance, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		var f domain.Fragrance
		if err := json.Unmarshal(hit.Source, &f); err != nil {
			return nil, fmt.Errorf("unmarshal hit: %w", err)
		}
		fragrances = append(fragrances, f)
	}

	facets := domain.Facets{
		Gender:      extractTermsAgg(result, "gender"),
		MainAccords: extractTermsAgg(result, "main_accords"),
	}

	return &domain.SearchResponse{
		Results: fragrances,
		Facets:  facets,
	}, nil
}

func extractTermsAgg(result *elastic.SearchResult, name string) map[string]int {
	buckets, found := result.Aggregations.Terms(name)
	if !found {
		return map[string]int{}
	}
	m := make(map[string]int, len(buckets.Buckets))
	for _, b := range buckets.Buckets {
		m[b.Key.(string)] = int(b.DocCount)
	}
	return m
}

func (a *ESFragranceRepositoryAdapter) SearchSimilar(ctx context.Context, id string, size int) ([]domain.Fragrance, error) {
	mlt := elastic.NewMoreLikeThisQuery().
		Field("name", "description", "main_accords").
		LikeItems(
			elastic.NewMoreLikeThisQueryItem().
				Index(indexName).
				Id(id),
		).
		MinTermFreq(1).
		MaxQueryTerms(12).
		MinDocFreq(2)

	q := elastic.NewBoolQuery().
		Must(mlt).
		MustNot(elastic.NewTermQuery("_id", id))

	result, err := a.db.Client.Search().
		Index(indexName).
		Query(q).
		Size(size).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("find similar: %w", err)
	}

	fragrances := make([]domain.Fragrance, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		var f domain.Fragrance
		if err := json.Unmarshal(hit.Source, &f); err != nil {
			return nil, fmt.Errorf("unmarshal hit: %w", err)
		}
		fragrances = append(fragrances, f)
	}
	return fragrances, nil
}
