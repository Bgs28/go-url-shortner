package model

import "time"

type URL struct {
	ID          int       `json:"id"`
	OriginalURL string    `json:"original_url"`
	ShortCode   string    `json:"short_code"`
	Clicks      int       `json:"clicks"`
	CreatedAt   time.Time `json:"created_at"`
}

type ShortenRequest struct {
	OriginalURL string `json:"original_url"`
}

type ShortenResponse struct {
	ShortCode 	string `json:"short_code"`
	ShortURL 	string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type StatsResponse struct {
	ShortCode 	string 	`json:"short_code"`
	OriginalURL string 	`json:"original_url"`
	Clicks      int       `json:"clicks"`
	CreatedAt   time.Time `json:"created_at"`
}