package fragrance

import (
	"context"
	"net/http"

	domain "github.com/artuos/sniffer/internal/domain/fragrance"
	schema "github.com/artuos/sniffer/internal/schema/fragrance"
	"github.com/gin-gonic/gin"
)

const errKey = "error"

type Service interface {
	Search(ctx context.Context, params schema.SearchParams) (*schema.SearchResponse, error)
	SearchSimilar(ctx context.Context, id string) ([]domain.Fragrance, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) *gin.RouterGroup {
	v1 := r.Group("/api/v1/fragrances")
	v1.GET("/search", h.Search)
	v1.GET("/:id/similar", h.SearchSimilar)
	return v1
}

func (h *Handler) Search(c *gin.Context) {
	params := schema.SearchParams{
		Query:  c.Query("q"),
		Gender: c.Query("gender"),
		Accord: c.Query("accord"),
	}

	if params.Query == "" && params.Gender == "" && params.Accord == "" {
		c.JSON(http.StatusBadRequest, gin.H{errKey: "at least one of 'q', 'gender', or 'accord' is required"})
		return
	}

	resp, err := h.service.Search(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{errKey: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) SearchSimilar(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{errKey: "path parameter 'id' is required"})
		return
	}

	results, err := h.service.SearchSimilar(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{errKey: err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
