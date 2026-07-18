package main

import (
	"context"
	"fmt"
	"os"

	"github.com/artuos/sniffer/internal/infra/adapters/extractor"
	"github.com/artuos/sniffer/internal/infra/adapters/repository/elasticsearch"
	schema "github.com/artuos/sniffer/internal/schema/fragrance"
	"github.com/artuos/sniffer/internal/service/fragrance"
)

func main() {
	ctx := context.Background()
	esIngestor, err := elasticsearch.NewESClientAdapter(ctx)
	if err != nil {
		panic(err)
	}
	defer esIngestor.Close()

	esIngestor.CreateIndex(ctx, "fragrances")

	fragranceRepo := elasticsearch.NewESFragranceRepositoryAdapter(esIngestor)
	ingestDataset := extractor.NewIngestDatasetAdapter[schema.FragranceModel]()
	fragranceService := fragrance.NewService(fragranceRepo, ingestDataset)

	fragranceService.IngestFragrances(ctx, os.Getenv("DATASET_PATH"))

	fmt.Printf("Dataset ingested successfully\n")
}
