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
	}
	return jsonRequest
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
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
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
