package ports

import "context"

type IngestDatasetPort[T any] interface {
	ConvertToModel(ctx context.Context, loc string) ([]T, error)
}
