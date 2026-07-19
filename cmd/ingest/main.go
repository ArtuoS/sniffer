package main

import (
	"context"
	"os"

	"github.com/artuos/sniffer/internal/config/container"
	"github.com/artuos/sniffer/internal/infra/adapters/downloader"
	"github.com/artuos/sniffer/internal/infra/adapters/extractor"
	"github.com/artuos/sniffer/internal/infra/adapters/repository/elasticsearch"
	schema "github.com/artuos/sniffer/internal/schema/fragrance"
	"github.com/artuos/sniffer/internal/service/fragrance"
	"go.uber.org/zap"
)

func main() {
	if err := container.Load(); err != nil {
		panic(err)
	}
	defer container.SyncLogger()
	logger := container.GetLogger()

	ctx := context.Background()
	esIngestor, err := elasticsearch.NewESClientAdapter(ctx)
	if err != nil {
		logger.Fatal("failed to create elasticsearch client", zap.Error(err))
	}
	defer esIngestor.Close()

	if err := esIngestor.CreateIndex(ctx, "fragrances"); err != nil {
		logger.Fatal("failed to create index", zap.Error(err))
	}

	fragranceRepo := elasticsearch.NewESFragranceRepositoryAdapter(esIngestor)
	ingestDataset := extractor.NewIngestDatasetAdapter[schema.FragranceModel]()
	kaggleDL := downloader.NewKaggleDownloader()
	fragranceService := fragrance.NewService(fragranceRepo, ingestDataset, kaggleDL)

	if err := fragranceService.IngestFragrancesFromKaggle(ctx, os.Getenv("KAGGLE_DATASET_URL")); err != nil {
		logger.Fatal("failed to ingest dataset", zap.Error(err))
	}

	logger.Info("dataset ingested successfully")
}
