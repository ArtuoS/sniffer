package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	httpHandler "github.com/artuos/sniffer/internal/http"
	"github.com/artuos/sniffer/internal/infra/adapters/repository/elasticsearch"
	"github.com/artuos/sniffer/internal/services"
)

func main() {
	ctx := context.Background()

	esClient, err := elasticsearch.NewESClientAdapter(ctx)
	if err != nil {
		log.Fatalf("connect to elasticsearch: %v", err)
	}
	defer esClient.Close()

	repo := elasticsearch.NewESFragranceRepositoryAdapter(esClient)
	svc := services.NewService(repo, nil)
	handler := httpHandler.NewHandler(svc)

	r := gin.Default()
	handler.RegisterRoutes(r)

	fmt.Println("API server running on :8081")
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("start server: %v", err)
	}
}
