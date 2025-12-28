package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mattnickolaus/chirpy/internal/database"

	"github.com/google/uuid"
)

var PROFANITY = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
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

func (cfg *apiConfig) getAllChrips(w http.ResponseWriter, r *http.Request) {
	allChirps, err := cfg.db.GetAllChrips(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error Reading Chirps from DB", err)
		return
	}

	responseChirps := []Chirp{}
	for _, c := range allChirps {
		newChirp := Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt.Time,
			UpdatedAt: c.UpdatedAt.Time,
			Body:      c.Body,
			UserID:    c.UserID,
		}
		responseChirps = append(responseChirps, newChirp)
	}

	respondWithJSON(w, http.StatusOK, responseChirps)
}

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	type chripRead struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	c := chripRead{}

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

	chirpParam := database.CreateChirpParams{
		Body:   cleanedChirp,
		UserID: c.UserID,
	}
	writenChirp, err := cfg.db.CreateChirp(r.Context(), chirpParam)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Chirp failed to write to database", err)
		return
	}

	returnedChirp := Chirp{
		ID:        writenChirp.ID,
		CreatedAt: writenChirp.CreatedAt.Time,
		UpdatedAt: writenChirp.UpdatedAt.Time,
		Body:      writenChirp.Body,
		UserID:    writenChirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, returnedChirp)
}
