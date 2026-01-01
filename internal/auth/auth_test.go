package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPassword(t *testing.T) {
	tests := []struct {
		name            string
		passwordSet     string
		passwordEntered string
		want            bool
	}{
		{
			name:            "Test 1: Standard Password Match",
			passwordSet:     "password123",
			passwordEntered: "password123",
			want:            true,
		},
		{
			name:            "Test 2: Standard Password Does Not Match",
			passwordSet:     "password",
			passwordEntered: "password123",
			want:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashed, err := HashPassword(tt.passwordSet)
			if err != nil {
				t.Errorf("HashPassword errored: %v", err)
			}

			actual, err := CheckPasswordHash(tt.passwordEntered, hashed)
			if err != nil {
				t.Errorf("CheckPasswordHash errored: %v", err)
			}

			expected := tt.want
			if actual != expected {
				t.Errorf("got: %v; want: %v", actual, expected)
			}

		})
	}
}

func TestValidateJWT(t *testing.T) {
	tests := []struct {
		name                  string
		userIDOnCreate        string
		tokenSecretOnCreate   string
		expiresIn             time.Duration
		tokenSecretOnValidate string
		match                 bool
	}{
		{
			name:                  "Test 1: JWT Correct Keys",
			userIDOnCreate:        "d9a6c7d5-de09-47c9-b8e0-d929e8af506c",
			tokenSecretOnCreate:   "password123",
			expiresIn:             time.Hour,
			tokenSecretOnValidate: "password123",
			match:                 true,
		},
		{
			name:                  "Test 2: JWT Non-Matching Keys",
			userIDOnCreate:        "d9a6c7d5-de09-47c9-b8e0-d929e8af506c",
			tokenSecretOnCreate:   "password123",
			expiresIn:             time.Hour,
			tokenSecretOnValidate: "password",
			match:                 false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := uuid.Parse(tt.userIDOnCreate)
			if err != nil {
				t.Errorf("Parsing userID string errored: %v\n", err)
			}

			tokenString, err := MakeJWT(userID, tt.tokenSecretOnCreate, tt.expiresIn)
			if err != nil {
				t.Errorf("MakeJWT errored: %v\n", err)
			}
			fmt.Printf("Token String: %v\n", tokenString)

			// Purposely ignoring error to test the outupt match
			actualUserID, _ := ValidateJWT(tokenString, tt.tokenSecretOnValidate)

			actualMatch := actualUserID == userID
			if actualMatch != tt.match {
				t.Errorf("got: %v; want: %v\nInputUser: %v; OutputUser: %v\n", actualMatch, tt.match, userID, actualUserID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name        string
		tokenString string
		match       bool
	}{
		{
			name:        "Test 1: Bearer Token String Passed",
			tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHkiLCJzdWIiOiJkOWE2YzdkNS1kZTA5LTQ3YzktYjhlMC1kOTI5ZThhZjUwNmMiLCJleHAiOjE3NjcyOTcxNTMsImlhdCI6MTc2NzI5MzU1M30.V-01iE9DyXoLuTc0qYyqVd3Ta7wiuhShfN0P7x_ILpU",
			match:       true,
		},
		{
			name:        "Test 2: Bearer Token with value Bearer",
			tokenString: "Bearer",
			match:       true,
		},
		{
			name:        "Test 3: Bearer Token with spaces in the value",
			tokenString: "Token String Value",
			match:       true,
		},
		{
			name:        "Test 4: Bearer Token Empty",
			tokenString: "",
			match:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := make(http.Header)
			authHeaderValue := fmt.Sprintf("Bearer %v", tt.tokenString)

			header.Add("Authorization", authHeaderValue)

			// Purposely ignoring error to test the outupt match
			outputTokenString, _ := GetBearerToken(header)

			fmt.Printf("Output Token String:'%v'\n", outputTokenString)

			actualMatch := outputTokenString == tt.tokenString
			if actualMatch != tt.match {
				t.Errorf("got: %v; want: %v\nInputUser: %v; OutputUser: %v\n", actualMatch, tt.match, outputTokenString, tt.tokenString)
			}
		})
	}
}
