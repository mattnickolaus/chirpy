package main

import (
	"testing"
)

func TestFilterProfanity(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Test 1: With Profanity",
			input: "This is a kerfuffle opinion I need to share with the world",
			want:  "This is a **** opinion I need to share with the world",
		},
		{
			name:  "Test 2: Normal",
			input: "Normal sentence without profanity",
			want:  "Normal sentence without profanity",
		},
		{
			name:  "Test 3: Profanity but upper",
			input: "I really need a kerfuffle to go to bed sooner, Fornax !",
			want:  "I really need a **** to go to bed sooner, **** !",
		},
		{
			name:  "Test 4: Profanity with puntuation",
			input: "I really need a kerfuffle to go to bed sooner, Fornax!",
			want:  "I really need a **** to go to bed sooner, Fornax!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := filterProfanity(tt.input)
			expected := tt.want

			if actual != expected {
				t.Errorf("got: %v; want: %v", actual, expected)
			}

		})
	}
}
