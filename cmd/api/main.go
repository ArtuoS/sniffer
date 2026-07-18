package main

import (
	"context"
	"fmt"
	"log"
	"os"

	fragranceHttp "github.com/artuos/sniffer/internal/http/fragrance"
	"github.com/artuos/sniffer/internal/infra/adapters/repository/elasticsearch"
	"github.com/artuos/sniffer/internal/service/fragrance"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	esClient, err := elasticsearch.NewESClientAdapter(ctx)
	if err != nil {
		log.Fatalf("connect to elasticsearch: %v", err)
	}
	defer esClient.Close()

	repo := elasticsearch.NewESFragranceRepositoryAdapter(esClient)
	svc := fragrance.NewService(repo, nil)
	handler := fragranceHttp.NewHandler(svc)

	r := gin.Default()
	handler.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	fmt.Printf("API server running on :%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("start server: %v", err)
	}
}

func registerRoutes(r *gin.Engine) {

}
