package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dmarab2/bot-request-site/backend/internal/database"
)

func turnRequestToJSON(databaseRequest database.Request) jsonRequest {
	jsonRequest := jsonRequest{
		ID:          int(databaseRequest.ID),
		CreatedAt:   databaseRequest.CreatedAt,
		UpdatedAt:   databaseRequest.UpdatedAt,
		RequestText: databaseRequest.RequestText,
		Status:      string(databaseRequest.Status),
	}
	return jsonRequest
}

func turnClaimToJson(databaseClaim database.RequestClaim) jsonClaim {
	jsonClaim := jsonClaim{
		RequestID: databaseClaim.RequestID,
		ClaimedAt: databaseClaim.ClaimedAt,
		ExpiresAt: databaseClaim.ExpiresAt,
	}
	return jsonClaim
}

// helper function to write a JSON error if something goes wrong during handling
func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorStruct struct {
		ErrorString string `json:"error"`
	}
	jsonErr := errorStruct{
		ErrorString: msg,
	}
	data, err := json.Marshal(jsonErr)
	if err != nil {
		log.Printf("Error marshaling json: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

// basic helper function to write a JSON bytearray to the address of the handler that called it
func respondWithJSON[T any](w http.ResponseWriter, code int, payload T) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling json: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

// metadataMiddleware is a function that adds metadata to any outgoing JSON request.
// this currently needs to be updated later once I figure out the necessary metadata to send.
func metadataMiddleware[T any](cfg *apiConfig, w http.ResponseWriter, code int, payload T) {
	type metadataPayload struct {
		Data      T    `json:"data"`
		PageNum   int  `json:"page_number"`
		NextLimit bool `json:"next_limit"`
		PrevLimit bool `json:"prev_limit"`
	}
	metaPay := metadataPayload{payload, 5, false, false}
	respondWithJSON(w, code, metaPay)
}
