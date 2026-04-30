package main

import (
	"context"

	"github.com/dmarab2/bot-request-site/backend/internal/database"
)

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
