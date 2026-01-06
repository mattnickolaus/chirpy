package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type polkaWebHookRequest struct {
	Event string `json:"event"`
	Data  struct {
		User_ID string `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) upgradeToChirpyRed(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	p := polkaWebHookRequest{}

	err := decoder.Decode(&p)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	if p.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	userID, err := uuid.Parse(p.Data.User_ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing given user ID", err)
		return
	}
	// We don't want write back the updated user, so ignoring the response
	_, err = cfg.db.UpgradeToChirpyRed(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error user could not be found when writing to database", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
