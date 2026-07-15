package handler

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/artuos/sniffer/internal/domain"
)

type Service interface {
	Search(ctx context.Context, query string) ([]domain.Fragrance, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	v1.GET("/fragrances/search", h.Search)
}
