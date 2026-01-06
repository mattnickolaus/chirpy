package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

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
		ID:          user.ID,
		CreatedAt:   user.CreatedAt.Time,
		UpdatedAt:   user.UpdatedAt.Time,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed.Bool,
	}

	respondWithJSON(w, http.StatusCreated, returnUser)
}

func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
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

	type updateUserInput struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	u := updateUserInput{}
	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&u)
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

	updateUserParams := database.UpdateUserParams{
		ID:             userID,
		Email:          u.Email,
		HashedPassword: convertedHashedPassword,
	}

	updatedUser, err := cfg.db.UpdateUser(r.Context(), updateUserParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error Writing to database", err)
		return
	}

	returnedUpdatedUser := User{
		ID:          updatedUser.ID,
		CreatedAt:   updatedUser.CreatedAt.Time,
		UpdatedAt:   updatedUser.UpdatedAt.Time,
		Email:       updatedUser.Email,
		IsChirpyRed: updatedUser.IsChirpyRed.Bool,
	}

	respondWithJSON(w, http.StatusOK, returnedUpdatedUser)
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

	// NOTE: Hard Coded 1 hour JWT expiration time
	expiresInHour := time.Second * time.Duration(3600)

	tokenString, err := auth.MakeJWT(user.ID, cfg.secret, expiresInHour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to generate Web Token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to generate Refresh Token", err)
		return
	}
	// NOTE: Hard Coded 60 days from now expiration for Refresh Token
	expiresSixtyDaysFromToday := sql.NullTime{
		Time:  time.Now().Add(60 * 24 * time.Hour),
		Valid: true,
	}
	refreshTokenParams := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		ExpiresAt: expiresSixtyDaysFromToday,
		UserID:    user.ID,
	}

	writtenRefreshToken, err := cfg.db.CreateRefreshToken(r.Context(), refreshTokenParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to write Refresh Token to DB", err)
		return
	}

	returnUser := User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
		Email:        user.Email,
		Token:        tokenString,
		RefreshToken: writtenRefreshToken.Token,
		IsChirpyRed:  user.IsChirpyRed.Bool,
	}

	respondWithJSON(w, http.StatusOK, returnUser)
}
