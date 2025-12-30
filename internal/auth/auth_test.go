package auth

import (
	"testing"
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
