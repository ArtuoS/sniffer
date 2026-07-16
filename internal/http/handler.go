package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/artuos/sniffer/internal/domain"
)

type Service interface {
	Search(ctx context.Context, params domain.SearchParams) (*domain.SearchResponse, error)
	SearchSimilar(ctx context.Context, id string) ([]domain.Fragrance, error)
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
	v1.GET("/fragrances/:id/similar", h.SearchSimilar)
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

func (h *Handler) SearchSimilar(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path parameter 'id' is required"})
		return
	}

	results, err := h.service.SearchSimilar(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
