package main

import (
	"context"
	"testing"
	"time"

	"github.com/dmarab2/bot-request-site/backend/internal/database"
)

// TestGetSingleRequestCore simply tests the logic in getSingleRequestCore to make sure it always returns a request according to
// the id it was given, even if the logic is updated and has middleware added to it. Note that the function core may change.
func TestGetSingleRequestCore(t *testing.T) {
	testContext := t.Context()
	testGetFunction := func(c context.Context, i int64) (database.Request, error) {
		testRequest := database.Request{
			ID:          i,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			RequestText: "TEST REQUEST",
			Status:      database.RequestStatusOpen}
		return testRequest, nil
	}
	var testID int64 = 5
	recievedRequest, err := getSingleRequestCore(testContext, testID, testGetFunction)
	if err != nil {
		t.Errorf("The test for getSingleRequestCore returned an error.")
	}
	if recievedRequest.ID != testID {
		t.Errorf("The test for getSingleRequestCore returned %d instead of %d.", recievedRequest.ID, testID)
	}

}

func testCreateRequestclaimCore(t *testing.T) {
}
