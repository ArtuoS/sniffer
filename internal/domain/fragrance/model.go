package fragrance

import "github.com/google/uuid"

type Fragrance struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Gender      string    `json:"gender"`
	RatingValue float64   `json:"rating_value"`
	RatingCount int64     `json:"rating_count"`
	MainAccords []string  `json:"main_accords"`
	Perfumers   string    `json:"perfumers"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
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
	Results []Fragrance `json:"results"`
	Facets  Facets      `json:"facets"`
}
