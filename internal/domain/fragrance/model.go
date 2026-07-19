package fragrance

import "github.com/google/uuid"

type Fragrance struct {
	Name        string    `json:"name"`
	Gender      string    `json:"gender"`
	Perfumers   string    `json:"perfumers"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	MainAccords []string  `json:"main_accords"`
	ID          uuid.UUID `json:"id"`
	RatingValue float64   `json:"rating_value"`
	RatingCount int64     `json:"rating_count"`
}

type SearchParams struct {
	Query  string
	Gender string
	Accord string
}

type Facets struct {
	Gender      map[string]int `json:"gender"`
	MainAccords map[string]int `json:"main_accords"`
}

type SearchResponse struct {
	Facets  Facets      `json:"facets"`
	Results []Fragrance `json:"results"`
}
