package extractor

import (
	"context"
	"fmt"
	"io"
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

func (i *IngestDatasetAdapter[T]) ConvertToModelFromReader(ctx context.Context, r io.Reader) ([]T, error) {
	var models []T
	if err := gocsv.Unmarshal(r, &models); err != nil {
		return nil, fmt.Errorf("unmarshal csv: %w", err)
	}
	return models, nil
}
