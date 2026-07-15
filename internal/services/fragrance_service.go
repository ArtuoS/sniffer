package services

import (
	"context"
	"fmt"

	"github.com/artuos/sniffer/internal/domain"
	"github.com/artuos/sniffer/internal/domain/ports/repository"
	"github.com/artuos/sniffer/internal/infra/ports"
	"github.com/artuos/sniffer/internal/schemas"
)

type Service struct {
	repository repository.FragranceRepositoryPort
	extractor  ports.IngestDatasetPort[schemas.FragranceModel]
}

func NewService(repository repository.FragranceRepositoryPort, extractor ports.IngestDatasetPort[schemas.FragranceModel]) *Service {
	return &Service{
		repository: repository,
		extractor:  extractor,
	}
}

func (s *Service) IngestFragrances(ctx context.Context, location string) error {
	models, err := s.extractor.ConvertToModel(ctx, location)
	if err != nil {
		fmt.Printf("Error converting dataset: %v\n", err)
		panic(err)
	}

	domainFragrances := make([]domain.Fragrance, len(models))
	for _, model := range models {
		fragrance := model.ToDomain()
		domainFragrances = append(domainFragrances, fragrance)
	}

	if err := s.repository.Create(ctx, domainFragrances); err != nil {
		return fmt.Errorf("failed to create fragrances: %w", err)
	}

	return nil
}

func (s *Service) Search(ctx context.Context, query string) ([]domain.Fragrance, error) {
	return s.repository.Search(ctx, query)
}
