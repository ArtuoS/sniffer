package fragrance

import (
	"context"
	"fmt"

	domain "github.com/artuos/sniffer/internal/domain/fragrance"
	repository "github.com/artuos/sniffer/internal/domain/fragrance/ports/repository"
	extractor "github.com/artuos/sniffer/internal/infra/ports/extractor"
	schema "github.com/artuos/sniffer/internal/schema/fragrance"
)

type Service struct {
	repo    repository.FragranceRepositoryPort
	extract extractor.IngestDatasetPort[schema.FragranceModel]
}

func NewService(repo repository.FragranceRepositoryPort, extract extractor.IngestDatasetPort[schema.FragranceModel]) *Service {
	return &Service{
		repo:    repo,
		extract: extract,
	}
}

func (s *Service) IngestFragrances(ctx context.Context, location string) error {
	models, err := s.extract.ConvertToModel(ctx, location)
	if err != nil {
		return fmt.Errorf("convert dataset: %w", err)
	}

	domainFragrances := make([]domain.Fragrance, 0, len(models))
	for _, model := range models {
		domainFragrances = append(domainFragrances, model.ToDomain())
	}

	if err := s.repo.Create(ctx, domainFragrances); err != nil {
		return fmt.Errorf("create fragrances: %w", err)
	}

	return nil
}

func (s *Service) Search(ctx context.Context, params domain.SearchParams) (*domain.SearchResponse, error) {
	return s.repo.Search(ctx, params)
}

func (s *Service) SearchSimilar(ctx context.Context, id string) ([]domain.Fragrance, error) {
	return s.repo.SearchSimilar(ctx, id, 10)
}
