// helper_functions mostly contains functions that are often used in the main package to write DRY code but aren't core server functionality.
// this includes JSON conversion, middleware, and response helpers.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"github.com/dmarab2/bot-request-site/backend/internal/database"
	"golang.org/x/text/unicode/norm"
)

var regexValidator string = `/[^a-z0-9()-_]/gi`

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

func turnLinkToJson(databaseTag database.RequestTag) jsonRequestTagLink {
	jsonLink := jsonRequestTagLink{
		RequestID: databaseTag.RequestID,
		TagID:     databaseTag.TagID,
	}
	return jsonLink
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

// normalizeTagName formats the names of any tags to follow a certain structure involving total lowercase and replacing spaces
// with underscores.
// for example "Ran Yakumo" will become "ran_yakumo."
func normalizeTagName(tagName string) string {
	// normalize Unicode in string
	tagName = norm.NFKC.String(tagName)
	// make every letter lowercase
	tagName = strings.ToLower(tagName)

	// replace spaces with underscores
	tagName = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return '_'
		}
		return r
	}, tagName)
	var buildString strings.Builder
	prevUnderscore := false

	// build a new string, remove repeated underscores
	for _, r := range tagName {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			buildString.WriteRune(r)
			prevUnderscore = false
		} else if r == '_' {
			if !prevUnderscore {
				buildString.WriteRune('_')
				prevUnderscore = true
			}
		}
	}
	// trim leading/trailing underscores
	finalString := strings.Trim(buildString.String(), "_")
	return finalString

}

func validateClaimPassword(newClaim requestClaimInsert) error {
	if newClaim.password == nil {
		return errors.New("User didn't provide a password.")
	}
	if *newClaim.password == "" {
		return errors.New("User submitted an empty password (how did they do that?)")
	}
	return nil
}

func createClaimParams(requestID int64, password string) database.CreateRequestClaimParams {
	claimParams := database.CreateRequestClaimParams{
		RequestID:       requestID,
		ClaimSecretHash: password,
	}
	return claimParams
}

func makeRequestParams(newInput ChangeStatusInput) database.ChangeRequestStatusParams {
	changedParams := database.ChangeRequestStatusParams{
		Status: database.RequestStatus(newInput.NewStatus),
		ID:     newInput.RequestID,
	}
	return changedParams
}

func validateChangeRequestStatus(input ChangeStatusInput) error {
	if input.RequestID <= 0 {
		return errors.New("Invalid request ID.")
	}
	switch input.NewStatus {
	case "open", "in_progress", "fulfilled", "cancelled":
		return nil
	default:
		return errors.New("Invalid request status.")
	}
}

func getTagByID(context context.Context, tagID int64, db *database.Queries) (database.Tag, error) {
	tag, err := db.GetTagByID(context, tagID)
	if err != nil {
		return database.Tag{}, errors.New("Could not get tag from database.")
	}
	return tag, nil
}

func validateTagLinkToRequest(input linkTagInput) error {
	if input.RequestID <= 0 {
		return errors.New("Invalid request ID.")
	}
	regexChecker, err := regexp.Compile(regexValidator)
	if err != nil {
		return errors.New("Unable to start the regex checking process.")
	}
	if regexChecker.MatchString(input.tagName) {
		return errors.New("Forbidden characters found in tag name!")
	}
	return nil
}

func makeTagLinkStruct(input linkTagInput) database.CreateRequestTagLinkParams {
	tagLinkParams := database.CreateRequestTagLinkParams{
		RequestID: input.RequestID,
		TagID:     input.tagID,
	}
	return tagLinkParams
}
