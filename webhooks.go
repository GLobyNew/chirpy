package main

import (
	"encoding/json"
	"net/http"

	"github.com/GLobyNew/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlePolkaWebhook(w http.ResponseWriter, r *http.Request) {
	apikey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error while getting api key")
		return
	}
	if apikey != cfg.polka_key {
		respondWithError(w, http.StatusUnauthorized, "invalid api key")
		return
	}

	type params struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	p := params{}
	err = decoder.Decode(&p)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error while decoding request")
		return
	}
	if p.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "event not supported")
		return
	}
	userUUID, err := uuid.Parse(p.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	_, err = cfg.db.UpgradeToChirpyRed(r.Context(), userUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "user not found")
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)

}
