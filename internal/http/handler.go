package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/artuos/sniffer/internal/domain"
)

type Service interface {
	Search(ctx context.Context, params domain.SearchParams) (*domain.SearchResponse, error)
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

func (h *Handler) Search(c *gin.Context) {
	params := domain.SearchParams{
		Query:  c.Query("q"),
		Gender: c.Query("gender"),
		Accord: c.Query("accord"),
	}

	if params.Query == "" && params.Gender == "" && params.Accord == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "at least one of 'q', 'gender', or 'accord' is required"})
		return
	}

	resp, err := h.service.Search(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
