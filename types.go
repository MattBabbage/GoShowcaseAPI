package main

import (
	"time"

	"github.com/google/uuid"
)

type CreateWeatherReportRequest struct {
	Description string  `json:"description"`
	Temperature int     `json:"temperature"`
	RainChance  float32 `json:"rain_chance"`
}

type WeatherReport struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description"`
	Temperature int       `json:"temperature"`
	RainChance  float32   `json:"rain_chance"`
}

func NewWeatherReport(Description string, Temperature int, RainChance float32) *WeatherReport {
	return &WeatherReport{
		ID:          uuid.New(),
		CreatedAt:   time.Now().UTC(),
		Description: Description,
		Temperature: Temperature,
		RainChance:  RainChance,
	}
}
