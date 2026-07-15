package adapters

import (
	"context"
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
)

type IngestDatasetAdapter[T any] struct {
}

func NewIngestDatasetAdapter[T any]() *IngestDatasetAdapter[T] {
	return &IngestDatasetAdapter[T]{}
}

func (i *IngestDatasetAdapter[T]) ConvertToModel(ctx context.Context, loc string) ([]T, error) {
	f, err := os.Open(loc)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var models []T
	if err := gocsv.UnmarshalFile(f, &models); err != nil {
		return nil, fmt.Errorf("unmarshal csv: %w", err)
	}

	return models, nil
}

// func convertToDocs(fragrances []FragranceModel) []map[string]any {
// 	docs := make([]map[string]any, 0, len(fragrances))
// 	for _, f := range fragrances {
// 		doc := map[string]any{
// 			"id": func() string {
// 				if f.ID != "" {
// 					return f.ID
// 				}
// 				return uuid.NewString()
// 			}(),
// 			"name":         f.Name,
// 			"gender":       f.Gender,
// 			"rating_value": float64(f.RatingValue),
// 			"rating_count": int64(f.RatingCount),
// 			"main_accords": []string(f.MainAccords),
// 			"perfumers":    f.Perfumers,
// 			"description":  f.Description,
// 			"url":          f.URL,
// 		}
// 		docs = append(docs, doc)
// 	}
// 	return docs
// }
