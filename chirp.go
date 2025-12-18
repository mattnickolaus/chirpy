package main

import (
	"encoding/json"
	"net/http"
)

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}
	type successResponse struct {
		Valid bool `json:"valid"`
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

	valid := successResponse{
		Valid: true,
	}
	respondWithJSON(w, http.StatusOK, valid)
}
