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

type jsonRequestTagLink struct {
	RequestID int64 `json:"request_id"`
	TagID     int64 `json:"tag_id"`
}

type requestClaimInsert struct {
	requestID int64
	password  *string
}

type ChangeStatusInput struct {
	RequestID int64
	NewStatus string
}

type linkTagInput struct {
	RequestID int64
	tagID     int64
	tagName   string
}
