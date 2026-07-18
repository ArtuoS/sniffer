package fragrance

import (
	"fmt"
	"strconv"
	"strings"

	domain "github.com/artuos/sniffer/internal/domain/fragrance"
	"github.com/google/uuid"
)

type FragranceModel struct {
	ID          string      `csv:"-"`
	Name        string      `csv:"Name"`
	Gender      string      `csv:"Gender"`
	RatingValue SafeFloat   `csv:"Rating Value"`
	RatingCount CommaInt64  `csv:"Rating Count"`
	MainAccords AccordsList `csv:"Main Accords"`
	Perfumers   string      `csv:"Perfumers"`
	Description string      `csv:"Description"`
	URL         string      `csv:"url"`
}

func parseOrCreateUUID(id string) uuid.UUID {
	if id == "" {
		return uuid.New()
	}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return uuid.New()
	}
	return parsedID
}

func (f *FragranceModel) ToDomain() domain.Fragrance {
	domainModel := domain.Fragrance{
		ID:          parseOrCreateUUID(f.ID),
		Name:        f.Name,
		Gender:      f.Gender,
		MainAccords: f.MainAccords,
		Perfumers:   f.Perfumers,
		Description: f.Description,
		URL:         f.URL,
	}
	if f.RatingValue != 0 {
		domainModel.RatingValue = float64(f.RatingValue)
	}
	if f.RatingCount != 0 {
		domainModel.RatingCount = int64(f.RatingCount)
	}
	return domainModel
}

type SafeFloat float64

func (f *SafeFloat) UnmarshalCSV(val string) error {
	if val == "" || val == "N/A" {
		return nil
	}
	v, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return fmt.Errorf("parse float %q: %w", val, err)
	}
	*f = SafeFloat(v)
	return nil
}

type CommaInt64 int64

func (c *CommaInt64) UnmarshalCSV(val string) error {
	if val == "" || val == "N/A" {
		return nil
	}
	cleaned := strings.ReplaceAll(val, ",", "")
	v, err := strconv.ParseInt(cleaned, 10, 64)
	if err != nil {
		return fmt.Errorf("parse int %q: %w", val, err)
	}
	*c = CommaInt64(v)
	return nil
}

type AccordsList []string

func (a *AccordsList) UnmarshalCSV(val string) error {
	val = strings.TrimSpace(val)
	val = strings.TrimPrefix(val, "[")
	val = strings.TrimSuffix(val, "]")

	if val == "" {
		return nil
	}

	parts := strings.Split(val, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		p = strings.Trim(p, "'")
		if p != "" {
			result = append(result, p)
		}
	}
	*a = result
	return nil
}
