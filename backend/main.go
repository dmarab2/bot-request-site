package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/dmarab2/bot-request-site/backend/internal/auth"
	"github.com/dmarab2/bot-request-site/backend/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// createRequestWriter is a function that takes an HTTP POST request to the API from the frontend and
// inserts a new request into the database. It runs on the "POST /api/requests" pattern.
func (cfg *apiConfig) createRequestWriter(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err.Error())
		w.WriteHeader(500)
		return
	}
	databaseRequest, err := cfg.db.CreateRequest(req.Context(), params.Body)
	if err != nil {
		log.Printf("Error inserting request: %s", err.Error())
		w.WriteHeader(500)
		return
	}
	jsonRequest := turnRequestToJSON(databaseRequest)
	respondWithJSON(w, 201, jsonRequest)
}

// getRequests returns a paginated list of requests in the database. It takes two query parameters: status, which
// corresponds to the status of the requests being asked for, and "after," which is ID of the request used for cursor pagination.
// Currently, a page is five requests, meaning only five requests are sent per query. This function is tied
// to the pattern "GET /api/requests"
func (cfg *apiConfig) getRequests(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Status database.RequestStatus
		ID     int64
	}
	requestStatus := req.URL.Query().Get("status")
	var requestSlice []database.Request
	var err error
	if requestStatus != "" {
		params := parameters{}
		params.Status = database.RequestStatus(requestStatus)
		cursorID := req.URL.Query().Get("after")
		if cursorID == "" {
			params.ID, err = cfg.db.GetFirstPageCursor(req.Context())
			if err != nil {
				log.Printf("Error getting request cursor: %s", err.Error())
				respondWithError(w, 500, "Could not get requests.")
				return
			}
			requestSlice, err = cfg.db.GetAllRequestsFiltered(req.Context(), database.GetAllRequestsFilteredParams(params))
			if err != nil {
				log.Printf("Error getting requests: %s\n", err.Error())
				respondWithError(w, 500, "Could not get requests.")
				return
			}
		} else {
			params.ID, err = strconv.ParseInt(req.URL.Query().Get("after"), 10, 64)
			if err != nil {
				log.Printf("Error parsing cursor: %s\n", err.Error())
				respondWithError(w, 500, "Could not get requests.")
				return
			}
			requestSlice, err = cfg.db.GetAllRequestsFiltered(req.Context(), database.GetAllRequestsFilteredParams(params))
			if err != nil {
				log.Printf("Error getting requests: %s\n", err.Error())
				respondWithError(w, 500, "Could not get requests.")
				return
			}
		}
	} else {
		requestSlice, err = cfg.db.GetAllRequests(req.Context())
		if err != nil {
			log.Printf("Error getting requests: %s\n", err.Error())
			respondWithError(w, 500, "Could not get all requests.")
			return
		}
	}
	jsonRequestSlice := make([]jsonRequest, 0, len(requestSlice))
	for _, request := range requestSlice {
		jsonRequest := turnRequestToJSON(request)
		jsonRequestSlice = append(jsonRequestSlice, jsonRequest)
	}
	log.Println(jsonRequestSlice)
	metadataMiddleware(cfg, w, 201, jsonRequestSlice)
}

// deleteRequests is a dev function to reset the database. DO NOT USE IN PROD. Tied to the
// "POST /admin/reset" pattern.
func (cfg *apiConfig) deleteRequests(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 404, "This is not the dev environment, you are not allowed to use this endpoint!")
	}
	err := cfg.db.DeleteAllRequests(req.Context())
	if err != nil {
		log.Printf("Error deleting requests: %s\n", err.Error())
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
}

// createRequestClaimWriter inserts a request claim into the database. The user submits a request ID and a password which
// ties the claim to a preexisting request and gives it a hashed password which can be checked for later updates.
func (cfg *apiConfig) createRequestClaimWriter(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		RequestID    int64  `json:"request_id"`
		PasswordHash string `json:"password"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error making a claim for this request: %s\n", err.Error())
		respondWithError(w, 500, "Error making claim")
		return
	}
	log.Printf("Raw password is: %s\n", params.PasswordHash)
	params.PasswordHash, err = auth.HashPassword(params.PasswordHash)
	if err != nil {
		log.Printf("Error making a claim for this request: %s\n", err.Error())
		respondWithError(w, 500, "Error making claim")
		return
	}
	log.Printf("Password hash is: %s\n", params.PasswordHash)
	databaseClaim, err := cfg.db.CreateRequestClaim(req.Context(), database.CreateRequestClaimParams{params.RequestID, params.PasswordHash})
	if err != nil {
		log.Printf("Error making a claim for this request: %s\n", err.Error())
		respondWithError(w, 500, "Error making claim")
		return
	}
	jsonClaim := turnClaimToJson(databaseClaim)
	respondWithJSON(w, 201, jsonClaim)
}

// main loads the .env variables, opens a connection to the postgres database, adds the endpoints the the server multiplexer
// and starts the server. Right now the server runs on port :8080.
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
	serveMux.HandleFunc("POST /admin/reset", cfg.deleteRequests)
	serveMux.HandleFunc("POST /api/requests", cfg.createRequestWriter)
	serveMux.HandleFunc("POST /api/request_claims", cfg.createRequestClaimWriter)
	serveMux.HandleFunc("GET /api/requests", cfg.getRequests)
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
		os.Exit(1)
	}
}
