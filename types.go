package main

import (
	"github.com/google/uuid"
)

type WeatherReport struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	Temperature int       `json:"temperature"`
	RainChance  float32   `json:"rain_chance"`
}

func NewWeatherReport(Description string, Temperature int, RainChance float32) *WeatherReport {
	return &WeatherReport{
		ID:          uuid.New(),
		Description: Description,
		Temperature: Temperature,
		RainChance:  RainChance,
	}
}
