package repository

import (
	"context"

	"github.com/artuos/sniffer/internal/domain"
)

type FragranceRepositoryPort interface {
	Create(ctx context.Context, models []domain.Fragrance) error
	Search(ctx context.Context, query string) ([]domain.Fragrance, error)
}
