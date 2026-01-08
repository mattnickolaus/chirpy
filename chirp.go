package main

import (
	"encoding/json"
	"net/http"
	"sort"
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
	authorID := r.URL.Query().Get("author_id")

	returnedChirps := []database.Chirp{}
	if authorID != "" {
		authorUserID, err := uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "No chirps by that author_id were found", err)
			return
		}
		returnedChirps, err = cfg.db.GetAllChirpsByUser(r.Context(), authorUserID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error: No chirps by that author_id were found", err)
			return
		}
	} else {
		allChirps, err := cfg.db.GetAllChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error Reading Chirps from DB", err)
			return
		}
		returnedChirps = allChirps
	}

	responseChirps := []Chirp{}
	for _, c := range returnedChirps {
		newChirp := Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt.Time,
			UpdatedAt: c.UpdatedAt.Time,
			Body:      c.Body,
			UserID:    c.UserID,
		}
		responseChirps = append(responseChirps, newChirp)
	}

	sortType := r.URL.Query().Get("sort")
	if sortType != "desc" && sortType != "asc" && sortType != "" {
		respondWithError(w, http.StatusNotFound, "Invalid sort parameter: accepts only 'asc' or 'desc'", nil)
		return
	}
	if sortType == "desc" {
		sort.Slice(responseChirps, func(i, j int) bool {
			return responseChirps[i].CreatedAt.After(responseChirps[j].CreatedAt)
		})
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

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
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

	queriedChirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp of that ID was not found in Database", err)
		return
	}
	if queriedChirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Not authorized to delete Chirp", nil)
		return
	}

	err = cfg.db.DeleteChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Chirp failed to detele to database", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
