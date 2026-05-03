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
