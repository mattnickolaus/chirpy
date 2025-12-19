package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

var PROFANITY = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}
	type successResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	c := chirp{}

	err := decoder.Decode(&c)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	chirpMaxLength := 140
	if len(c.Body) > chirpMaxLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanedChirp := filterProfanity(c.Body)

	cleanedResponse := successResponse{
		CleanedBody: cleanedChirp,
	}
	respondWithJSON(w, http.StatusOK, cleanedResponse)
}

func filterProfanity(chirpText string) string {

	chirpWords := strings.Split(chirpText, " ")

	for i, w := range chirpWords {
		lowercaseWord := strings.ToLower(w)
		if _, containsProf := PROFANITY[lowercaseWord]; containsProf {
			chirpWords[i] = "****"
		}
	}

	return strings.Join(chirpWords, " ")
}
