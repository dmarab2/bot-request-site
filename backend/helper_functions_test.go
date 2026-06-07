package main

import (
	"testing"
)

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
