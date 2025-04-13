package main

import (
	"log"
	"net/http"

	"github.com/GLobyNew/chirpy/internal/auth"
)

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, r *http.Request) {

	type resp struct {
		Token string `json:"token"`
	}

	if len(r.Header) == 0 {
		log.Println("no headers provided")
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	authToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Println("auth token is not present")
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	refreshTokenInDB, err := cfg.db.GetRefreshTokenByToken(r.Context(), authToken)
	if err != nil {
		log.Println("no token found in db")
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	cfg.db.RevokeToken(r.Context(), refreshTokenInDB.Token)

	respondWithJSON(w, http.StatusNoContent, nil)
}
