package main

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dmarab2/bot-request-site/backend/internal/auth"
	"github.com/dmarab2/bot-request-site/backend/internal/database"
)

// makes sure that turnRequestToJson accurately transfers request info over.
func TestTurnRequestToJson(t *testing.T) {
	testRequest := database.Request{
		// 25 is an arbitrary number
		ID:          int64(25),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		RequestText: "TestRequest",
		Status:      "fulfilled",
	}
	testJsonRequest := turnRequestToJSON(testRequest)
	if int(testRequest.ID) != testJsonRequest.ID ||
		testRequest.RequestText != testJsonRequest.RequestText ||
		string(testRequest.Status) != testJsonRequest.Status {
		t.Errorf("testRequest did not equate to testJsonStatus. Values should be %d, %s, and %s, but are %d, %s, and %s", testRequest.ID, testRequest.RequestText, string(testRequest.Status), testJsonRequest.ID, testJsonRequest.RequestText, testJsonRequest.Status)
	}
}

// makes sure that turnClaimToJson accurately transfers claim info over.
func TestTurnClaimToJson(t *testing.T) {
	testClaim := database.RequestClaim{
		// 25 is an arbitrary number
		RequestID:       int64(25),
		ClaimedAt:       time.Now(),
		ClaimSecretHash: "secrethash",
		ExpiresAt:       time.Now().AddDate(0, 0, 1),
	}
	testPassword := "testPassword"
	testJsonClaim := turnClaimToJson(testClaim, &testPassword)
	if int(testClaim.RequestID) != int(testJsonClaim.RequestID) {
		t.Errorf("testClaim ID %d did not equal testJsonClaim ID %d", testClaim.RequestID, testJsonClaim.RequestID)
	}
}

func TestTurnLinkToJson(t *testing.T) {
	testRequestLink := database.RequestTag{
		RequestID: int64(30),
		TagID:     int64(30),
	}
	testJsonLink := turnLinkToJson(testRequestLink)
	if testRequestLink.RequestID != testJsonLink.RequestID || testRequestLink.TagID != testJsonLink.TagID {
		t.Errorf("testRequestLink values were not transferred over to testJsonLink")
	}
}

// TestRespondWithError ensures that respondWithError properly writes the error passed to it.
func TestRespondWithError(t *testing.T) {
	type jsonDecoder struct {
		TestErrorHolder string `json:"error"`
	}
	w := httptest.NewRecorder()
	respondWithError(w, 500, "TEST ERROR")
	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != 500 {
		t.Errorf("Status code returned was %d instead of %d", resp.StatusCode, 500)
	}
	testStruct := jsonDecoder{}
	if err := json.NewDecoder(resp.Body).Decode(&testStruct); err != nil {
		t.Errorf("Error decoding the test error message: %s", err.Error())
	}
	if testStruct.TestErrorHolder != "TEST ERROR" {
		t.Errorf("testErrorHolder is %s instead of TEST ERROR", testStruct.TestErrorHolder)
	}

}

// TestRespondWithJson ensure that respondWithJson returns the code and data that has been passed to it.
func TestRespondWithJson(t *testing.T) {
	type testObject struct {
		TestStringProperty string `json:"test_string"`
	}
	w := httptest.NewRecorder()
	testJsonObject := testObject{TestStringProperty: "TEST PROPERTY"}
	respondWithJSON(w, 201, testJsonObject)
	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		t.Errorf("Status code returned was %d instead of %d", resp.StatusCode, 201)
	}
	receiverStruct := testObject{}
	if err := json.NewDecoder(resp.Body).Decode(&receiverStruct); err != nil {
		t.Errorf("Error decoding the test message: %s", err.Error())
	}
	if receiverStruct.TestStringProperty != "TEST PROPERTY" {
		t.Errorf("TestStringProperty is %s instead of TEST PROPERTY", receiverStruct.TestStringProperty)
	}

}

// TestValidateClaimPassword tests every possible output from validateClaimPassword to make sure that empty strings,
// nonexistent submissions, and valid submissions receive the correct response.
func TestValidateClaimPassword(t *testing.T) {
	testPassword := "testPassword"
	testClaim := requestClaimInsert{requestID: 40, password: &testPassword}
	err := validateClaimPassword(testClaim)
	if err != nil {
		t.Errorf("Test for validating claim password threw an error when it should've passed")
	}
	testPassword = ""
	testClaim.password = &testPassword
	err = validateClaimPassword(testClaim)
	if err == nil {
		t.Errorf("Test for validate claim password should have failed due to an empty string")
	}
	testClaim.password = nil
	err = validateClaimPassword(testClaim)
	if err == nil {
		t.Errorf("Test for validate claim password should have failed due to a nil pointer")
	}
}

func TestCreateClaimParams(t *testing.T) {
	testRequestID := int64(30)
	testPasswordHash := "testHash"
	testClaimParams := createClaimParams(testRequestID, testPasswordHash)
	if testClaimParams.RequestID != testRequestID || testClaimParams.ClaimSecretHash != testPasswordHash {
		t.Errorf("Claim params do not match the Request ID and Password Hash.")
	}
}

func TestMakeRequestParams(t *testing.T) {
	testNewInput := ChangeStatusInput{RequestID: 30, NewStatus: "in_progress"}
	testChangeParams := makeRequestParams(testNewInput)
	if testNewInput.RequestID != testChangeParams.ID || testNewInput.NewStatus != string(testChangeParams.Status) {
		t.Errorf("Input does not match the param changing object.")
	}
}

// TestValidateChangeRequestStatus ensures that new status changes are valid statuses.
func TestValidateChangeRequestStatus(t *testing.T) {
	testInput := ChangeStatusInput{20, "open"}
	err := validateChangeRequestStatus(testInput)
	if err != nil {
		t.Errorf("Test for validating request status failed with %s", err.Error())
	}
	testInput = ChangeStatusInput{-5, "open"}
	err = validateChangeRequestStatus(testInput)
	if err == nil {
		t.Errorf("Test for validating request status passed when it should've failed")
	}
	testInput = ChangeStatusInput{25, "BAD STATUS"}
	err = validateChangeRequestStatus(testInput)
	if err == nil {
		t.Errorf("Test for validating request status passed when it should've failed")
	}
}

// TestValidateTagLinkToRequest makes sure that new link tags are valid under several constraints.
func TestValidateTagLinkToRequest(t *testing.T) {
	testInput := linkTagInput{10, 10, "test"}
	err := validateTagLinkToRequest(testInput)
	if err != nil {
		t.Errorf("Test for validating tag links failed with %s", err.Error())
	}
	testInput.RequestID = -5
	err = validateTagLinkToRequest(testInput)
	if err == nil {
		t.Errorf("Tag link validation passed when it should've failed due to negative request ID")
	}
	testInput.RequestID = 10
	testInput.tagName = "&&$$@#*$"
	err = validateTagLinkToRequest(testInput)
	if err == nil {
		t.Errorf("Tag link validation passed when it should've failed due to invalid tag name")
	}
}

// TestNormalizeTagName makes sure that all strings are normalized to the same format no matter what the string contains.
func TestNormalizeTagName(t *testing.T) {
	testString := "Ran Yakumo"
	returnedString := normalizeTagName(testString)
	if returnedString != "ran_yakumo" {
		t.Errorf("String 'Ran Yakumo' should have converted to 'ran_yakumo' and didn't")
	}
	testString = "Test Space Underscore"
	returnedString = normalizeTagName(testString)
	if returnedString != "test_space_underscore" {
		t.Errorf("String 'Test Space Underscore' should have converted to 'test_space_underscore' and didn't")
	}
}

func TestPasswordGenAndHashing(t *testing.T) {
	testPassword, err := auth.GenerateClaimPassword()
	if err != nil {
		t.Errorf("Something went wrong when trying to gen a test claim password: %s", err.Error())
	}
	hashedTestPassword, err := auth.HashPassword(testPassword)
	if err != nil {
		t.Errorf("Something went wrong when trying to gen a test claim password: %s", err.Error())
	}
	doPasswordsMatch, err := auth.CheckPasswordHash(testPassword, hashedTestPassword)
	if err != nil {
		t.Errorf("Something went wrong when trying to gen a test claim password: %s", err.Error())
	}
	t.Logf("Password: %s, Hash: %s", testPassword, hashedTestPassword)
	if !doPasswordsMatch {
		t.Errorf("The password does not match the hash!")
	}
}
