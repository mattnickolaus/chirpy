package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type userEmail struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	u := userEmail{}

	err := decoder.Decode(&u)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), u.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error Writing to database", err)
		return
	}

	returnUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Email:     user.Email,
	}

	respondWithJSON(w, http.StatusCreated, returnUser)
}
