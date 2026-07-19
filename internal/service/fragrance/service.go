package fragrance

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/artuos/sniffer/internal/config/container"
	domain "github.com/artuos/sniffer/internal/domain/fragrance"
	repository "github.com/artuos/sniffer/internal/domain/fragrance/ports/repository"
	schema "github.com/artuos/sniffer/internal/schema/fragrance"
	"go.uber.org/zap"
)

type DatasetDownloader interface {
	FetchCSV(ctx context.Context, datasetURL string) (*bytes.Reader, error)
}

type DatasetConverter[T any] interface {
	ConvertToModel(ctx context.Context, loc string) ([]T, error)
	ConvertToModelFromReader(ctx context.Context, r io.Reader) ([]T, error)
}

type Service struct {
	repo       repository.FragranceRepositoryPort
	extract    DatasetConverter[schema.FragranceModel]
	downloader DatasetDownloader
	logger     *zap.Logger
}

func NewService(repo repository.FragranceRepositoryPort, extract DatasetConverter[schema.FragranceModel], downloader DatasetDownloader) *Service {
	return &Service{
		repo:       repo,
		extract:    extract,
		downloader: downloader,
		logger:     container.GetLogger(),
	}
}

func (s *Service) IngestFragrances(ctx context.Context, location string) error {
	models, err := s.extract.ConvertToModel(ctx, location)
	if err != nil {
		return fmt.Errorf("convert dataset: %w", err)
	}

	domainFragrances := make([]domain.Fragrance, 0, len(models))
	for i := range models {
		domainFragrances = append(domainFragrances, models[i].ToDomain())
	}

	if err := s.repo.Create(ctx, domainFragrances); err != nil {
		return fmt.Errorf("create fragrances: %w", err)
	}

	return nil
}

func (s *Service) IngestFragrancesFromKaggle(ctx context.Context, datasetURL string) error {
	csvReader, err := s.downloader.FetchCSV(ctx, datasetURL)
	if err != nil {
		return fmt.Errorf("download dataset: %w", err)
	}

	models, err := s.extract.ConvertToModelFromReader(ctx, csvReader)
	if err != nil {
		return fmt.Errorf("convert dataset: %w", err)
	}

	domainFragrances := make([]domain.Fragrance, 0, len(models))
	for i := range models {
		domainFragrances = append(domainFragrances, models[i].ToDomain())
	}

	if err := s.repo.Create(ctx, domainFragrances); err != nil {
		return fmt.Errorf("create fragrances: %w", err)
	}

	return nil
}

func (s *Service) Search(ctx context.Context, params schema.SearchParams) (*schema.SearchResponse, error) {
	return s.repo.Search(ctx, params)
}

func (s *Service) SearchSimilar(ctx context.Context, id string) ([]domain.Fragrance, error) {
	return s.repo.SearchSimilar(ctx, id, 10)
}
