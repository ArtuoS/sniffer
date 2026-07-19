package main

import (
	"context"
	"flag"
	"os"

	"github.com/artuos/sniffer/internal/config/container"
	fragranceHttp "github.com/artuos/sniffer/internal/http/fragrance"
	"github.com/artuos/sniffer/internal/infra/adapters/downloader"
	"github.com/artuos/sniffer/internal/infra/adapters/extractor"
	"github.com/artuos/sniffer/internal/infra/adapters/repository/elasticsearch"
	schema "github.com/artuos/sniffer/internal/schema/fragrance"
	"github.com/artuos/sniffer/internal/service/fragrance"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	mode := flag.String("mode", "api", "run mode: api or ingest")
	flag.Parse()

	if err := container.Load(); err != nil {
		panic(err)
	}
	defer container.SyncLogger()
	logger := container.GetLogger()

	switch *mode {
	case "api":
		runAPI(logger)
	case "ingest":
		runIngest(logger)
	default:
		logger.Fatal("unknown mode", zap.String("mode", *mode))
	}
}

func runAPI(logger *zap.Logger) {
	ctx := context.Background()

	esClient, err := elasticsearch.NewESClientAdapter(ctx)
	if err != nil {
		logger.Fatal("failed to connect to elasticsearch", zap.Error(err))
	}
	defer esClient.Close()

	repo := elasticsearch.NewESFragranceRepositoryAdapter(esClient)
	svc := fragrance.NewService(repo, nil, nil)
	handler := fragranceHttp.NewHandler(svc)

	r := gin.Default()
	handler.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	logger.Info("server started", zap.String("port", port))
	if err := r.Run(":" + port); err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}
}

func runIngest(logger *zap.Logger) {
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
