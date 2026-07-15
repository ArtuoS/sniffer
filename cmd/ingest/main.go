package main

import (
	"context"
	"fmt"

	"github.com/artuos/sniffer/internal/infra/adapters"
	"github.com/artuos/sniffer/internal/infra/adapters/repository/elasticsearch"
	"github.com/artuos/sniffer/internal/schemas"
	"github.com/artuos/sniffer/internal/services"
)

func main() {
	ctx := context.Background()
	esIngestor, err := elasticsearch.NewESClientAdapter(ctx)
	if err != nil {
		panic(err)
	}
	// defer esIngestor.Close()

	esIngestor.CreateIndex(ctx, "fragrances")

	fragranceRepo := elasticsearch.NewESFragranceRepositoryAdapter(esIngestor)
	ingestDataset := adapters.NewIngestDatasetAdapter[schemas.FragranceModel]()
	fragranceService := services.NewService(fragranceRepo, ingestDataset)

	fragranceService.IngestFragrances(ctx, "/mnt/d/Distros/Ubuntu/Programs/Projects/sniffer/dataset/fra_perfumes.csv")

	fmt.Printf("Dataset ingested successfully\n")

}
