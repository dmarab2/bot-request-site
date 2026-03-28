package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dmarab2/bot-request-site/backend/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type jsonRequest struct {
	ID          int       `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	RequestText string    `json:"request_text"`
}

type apiConfig struct {
	db     *database.Queries
	secret string
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

func turnRequestToJSON(databaseRequest database.Request) jsonRequest {
	jsonRequest := jsonRequest{
		ID:          int(databaseRequest.ID),
		CreatedAt:   databaseRequest.CreatedAt,
		UpdatedAt:   databaseRequest.UpdatedAt,
		RequestText: databaseRequest.RequestText,
	}
	return jsonRequest
}

func (cfg *apiConfig) createRequestWriter(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	databaseRequest, err := cfg.db.CreateRequest(req.Context(), params.Body)
	if err != nil {
		log.Printf("Error inserting request: %s", err)
		w.WriteHeader(500)
		return
	}
	jsonRequest := turnRequestToJSON(databaseRequest)
	respondWithJSON(w, 201, jsonRequest)
}

// A simple beginning function for now.
func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)
	cfg := &apiConfig{dbQueries, os.Getenv("SECRET_KEY")}
	serveMux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}
	serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested %s\n", r.URL.Path)
	})
	fileServer := http.FileServer(http.Dir("./static"))
	serveMux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	serveMux.HandleFunc("POST /api/requests", cfg.createRequestWriter)
	/*
		ticker := time.NewTicker(1 * time.Minute)
				go func() {
			        // Optional: Clean up if the main function ever returns
			        defer ticker.Stop()

			        for {
			            select {
			            case t := <-ticker.C:
			                // 3. Your periodic logic here
			                fmt.Println("Running background task at:", t)

			                // Example: You can access your cfg or db here
			                // cfg.db.SomeCleanupMethod(context.Background())
			            }
			        }
			    }()
	*/
	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("There was an error: %s\n", err.Error())
	}
}
