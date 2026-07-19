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
