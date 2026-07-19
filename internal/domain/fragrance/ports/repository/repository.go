package repository

import (
	"context"

	domain "github.com/artuos/sniffer/internal/domain/fragrance"
	schema "github.com/artuos/sniffer/internal/schema/fragrance"
)

type FragranceRepositoryPort interface {
	Create(ctx context.Context, models []domain.Fragrance) error
	Search(ctx context.Context, params schema.SearchParams) (*schema.SearchResponse, error)
	SearchSimilar(ctx context.Context, id string, size int) ([]domain.Fragrance, error)
}
