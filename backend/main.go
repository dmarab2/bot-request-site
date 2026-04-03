package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dmarab2/bot-request-site/backend/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// function to create a single request. only takes the request text itself, the other fields are propagated automatically.
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

func (cfg *apiConfig) deleteRequests(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 404, "This is not the dev environment, you are not allowed to use this endpoint!")
	}
	err := cfg.db.DeleteAllRequests(req.Context())
	if err != nil {
		log.Printf("Error deleting requests: %s\n", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
}

// A simple beginning function for now.
func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	fmt.Printf("dburl is %s\n", dbURL)
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)
	cfg := &apiConfig{dbQueries, os.Getenv("SECRET_KEY"), os.Getenv("PLATFORM")}
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
