// core_functions.go contains the core, "pure" logic of the handler functions in main.go.
// I separated the logic for a few reasons, one being to make them easier to test, another to prevent changes to the way I grab
// the data from the frontend from affecting the core logic, and also because I needed practice with writing code in this manner.
// Every function thus takes in dependencies and does not mutate anything outside of itself.
package main

import (
	"context"

	"github.com/dmarab2/bot-request-site/backend/internal/database"
)

// changeRequestStatusCore is the core logic of the changeRequestStatus writer function in main. It checks to make sure the new status
// passed in is valid, then it makes a parameters struct from it and adds the change to the database using a passed in dependency.
func changeRequestStatusCore(
	context context.Context,
	input ChangeStatusInput,
	updateFunction func(context.Context, database.ChangeRequestStatusParams) (database.Request, error),
) (database.Request, error) {
	if err := validateChangeRequestStatus(input); err != nil {
		return database.Request{}, err
	}
	params := makeRequestParams(input)
	return updateFunction(context, params)
}

// createRequestClaimCore is the core logic of the createRequestClaimWriter function in main. It validates the provided password,
// creates a parameters struct from it, and then passes it into the dependency function to insert it into the database.
func createRequestClaimCore(
	context context.Context,
	newClaim requestClaimInsert,
	hashedPassword string,
	insertFunction func(context.Context, database.CreateRequestClaimParams) (database.RequestClaim, error),
) (database.RequestClaim, error) {
	if err := validateClaimPassword(newClaim); err != nil {
		return database.RequestClaim{}, err
	}
	claimParams := createClaimParams(newClaim.requestID, hashedPassword)
	return insertFunction(context, claimParams)
}

// linkTagToRequestCore is the core logic of the linkTagToRequest handler in main. It creates a tagLink input struct that is validated
// by checking to see if the id and name is valid, and then uses the linker function to add a row to the linker table.
func linkTagToRequestCore(
	context context.Context,
	requestID int64,
	relevantTag database.Tag,
	linkerFunction func(context.Context, database.CreateRequestTagLinkParams) (database.RequestTag, error),
) (database.RequestTag, error) {
	tagLink := makeTagLinkInput(requestID, relevantTag)
	if err := validateTagLinkToRequest(tagLink); err != nil {
		return database.RequestTag{}, err
	}
	tagLinkParams := makeTagLinkParams(tagLink)
	return linkerFunction(context, tagLinkParams)
}

// getSingleRequestCore is a simple function, and simply returns the result of pulling a single request from the database.
func getSingleRequestCore(
	context context.Context,
	requestID int64,
	getFunction func(context.Context, int64) (database.Request, error),
) (database.Request, error) {
	return getFunction(context, requestID)
}
