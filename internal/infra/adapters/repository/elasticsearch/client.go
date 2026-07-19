package elasticsearch

import (
	"context"
	"fmt"
	"os"

	"github.com/artuos/sniffer/internal/config/container"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
)

type ESClientAdapter struct {
	Client *elastic.Client
	logger *zap.Logger
}

func NewESClientAdapter(ctx context.Context) (*ESClientAdapter, error) {
	esURL := os.Getenv("ELASTICSEARCH_URL")
	client, err := elastic.NewSimpleClient(
		elastic.SetURL(esURL),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(true),
	)
	if err != nil {
		return nil, fmt.Errorf("create elastic client: %w", err)
	}

	info, code, err := client.Ping(esURL).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping elasticsearch: %w", err)
	}

	logger.Info("elasticsearch connected",
		zap.Int("status_code", code),
		zap.String("version", info.Version.Number),
	)

	return &ESClientAdapter{
		Client: client,
		logger: container.GetLogger(),
	}, nil
}

func (a *ESClientAdapter) CreateIndex(ctx context.Context, indexName string) error {
	exists, err := a.Client.IndexExists(indexName).Do(ctx)
	if err != nil {
		return fmt.Errorf("check index exists: %w", err)
	}
	if exists {
		if _, err := a.Client.DeleteIndex(indexName).Do(ctx); err != nil {
			return fmt.Errorf("delete existing index: %w", err)
		}
	}

	if _, err = a.Client.CreateIndex(indexName).Body(fragranceMapping).Do(ctx); err != nil {
		return fmt.Errorf("create index: %w", err)
	}
	return nil
}

func (a *ESClientAdapter) Close() {
	if a.Client != nil {
		a.Client.Stop()
	}
}
