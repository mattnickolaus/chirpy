package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mattnickolaus/chirpy/internal/auth"
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

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDPath := r.PathValue("chirpID")
	if chirpIDPath == "" {
		respondWithError(w, http.StatusBadRequest, "Unable to retrieve chirpID from path", nil)
		return
	}
	chirpID, err := uuid.Parse(chirpIDPath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to parse the uuid from provide path", err)
		return
	}

	readChirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp was not read from database", err)
		return
	}

	returnedChirp := Chirp{
		ID:        readChirp.ID,
		CreatedAt: readChirp.CreatedAt.Time,
		UpdatedAt: readChirp.UpdatedAt.Time,
		Body:      readChirp.Body,
		UserID:    readChirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, returnedChirp)
}

func (cfg *apiConfig) getAllChirps(w http.ResponseWriter, r *http.Request) {
	allChirps, err := cfg.db.GetAllChirps(r.Context())
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
		Body string `json:"body"`
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized: Token Invalid", err)
		return
	}
	userID, err := auth.ValidateJWT(tokenString, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized: Token Invalid:", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	c := chripRead{}

	err = decoder.Decode(&c)
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
		UserID: userID,
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
