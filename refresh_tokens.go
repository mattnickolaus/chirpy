package main

import (
	"net/http"
	"time"

	"github.com/mattnickolaus/chirpy/internal/auth"
)

func (cfg *apiConfig) refreshAccessToken(w http.ResponseWriter, r *http.Request) {
	refreshTokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized: Refresh Token Invalid", err)
		return
	}

	refreshTokenRecord, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshTokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized: Refresh Token Invalid", err)
		return
	}

	// NOTE: Hard Coded 1 hour JWT expiration time
	expiresInHour := time.Second * time.Duration(3600)
	accessTokenString, err := auth.MakeJWT(refreshTokenRecord.UserID, cfg.secret, expiresInHour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to generate Web Token", err)
		return
	}

	type accessTokenResponse struct {
		Token string `json:"token"`
	}

	returnedAccessToken := accessTokenResponse{
		Token: accessTokenString,
	}

	respondWithJSON(w, http.StatusOK, returnedAccessToken)
}
