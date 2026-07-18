package elasticsearch

import "context"

type ESClientPort interface {
	CreateIndex(ctx context.Context) error
	IndexFragrances(ctx context.Context, docs []map[string]any) error
}
