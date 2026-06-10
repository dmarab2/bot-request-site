package main

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

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

// TestNormalizeTagName makes sure that all strings are normalized to the same format no matter what the string contains.
func TestNormalizeTagName(t *testing.T) {
	testString := "Ran Yakumo"
	returnedString := normalizeTagName(testString)
	if returnedString != "ran_yakumo" {
		t.Errorf("String 'Ran Yakumo' should have converted to 'ran_yakumo' and didn't")
	}
}
