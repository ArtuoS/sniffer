package repository

import (
	"context"

	domain "github.com/artuos/sniffer/internal/domain/fragrance"
)

type FragranceRepositoryPort interface {
	Create(ctx context.Context, models []domain.Fragrance) error
	Search(ctx context.Context, params domain.SearchParams) (*domain.SearchResponse, error)
	SearchSimilar(ctx context.Context, id string, size int) ([]domain.Fragrance, error)
}
