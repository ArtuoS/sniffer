package main

import (
	"context"
	"os"

	"github.com/artuos/sniffer/internal/config/container"
	fragranceHttp "github.com/artuos/sniffer/internal/http/fragrance"
	"github.com/artuos/sniffer/internal/infra/adapters/repository/elasticsearch"
	"github.com/artuos/sniffer/internal/service/fragrance"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	if err := container.Load(); err != nil {
		panic(err)
	}
	defer container.SyncLogger()
	logger := container.GetLogger()
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
