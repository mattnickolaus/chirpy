package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/mattnickolaus/chirpy/internal/auth"
	"github.com/mattnickolaus/chirpy/internal/database"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type userInput struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	u := userInput{}

	err := decoder.Decode(&u)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	hashedPassword, err := auth.HashPassword(u.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}
	convertedHashedPassword := sql.NullString{String: hashedPassword, Valid: true}

	userParams := database.CreateUserParams{
		Email:          u.Email,
		HashedPassword: convertedHashedPassword,
	}

	user, err := cfg.db.CreateUser(r.Context(), userParams)
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

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	type loginInput struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	u := loginInput{}

	err := decoder.Decode(&u)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	user, err := cfg.db.GetUserByUsername(r.Context(), u.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	matched, err := auth.CheckPasswordHash(u.Password, user.HashedPassword.String)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	if !matched {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	returnUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Email:     user.Email,
	}

	respondWithJSON(w, http.StatusOK, returnUser)
}
