package domain

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
