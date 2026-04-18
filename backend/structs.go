// structs.go contains helper structs that are used in the main package. These include json versions of requests and claims to
// return in an HTTP request to a user, and the apiConfig object that holds links to the database.
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

type jsonClaim struct {
	RequestID int64     `json:"request_id"`
	ClaimedAt time.Time `json:"claimed_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
