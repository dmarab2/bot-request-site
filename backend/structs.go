package main

import (
	"time"

	"github.com/dmarab2/bot-request-site/backend/internal/database"
)

type jsonRequest struct {
	ID          int       `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	RequestText string    `json:"request_text"`
	Status      string    `json:"status_"`
}

type apiConfig struct {
	db       *database.Queries
	secret   string
	platform string
}
