package main

import (
	"log"
	"net/http"
	"time"

	"github.com/GLobyNew/chirpy/internal/auth"
)

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, r *http.Request) {
	log.SetPrefix("handleRefresh: ")
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

	tokenInDB, err := cfg.db.GetRefreshTokenByToken(r.Context(), authToken)
	if err != nil {
		log.Println("no token found in db")
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	if tokenInDB.ExpiresAt.Before(time.Now()) {
		log.Println("token is expired")
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	if tokenInDB.RevokedAt.Valid {
		log.Println("token is revoked")
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	// Generate JWT token
	jwtToken, err := auth.MakeJWT(tokenInDB.UserID, cfg.jwtSecret, DefaultExpiresIn)
	if err != nil {
		log.Println("error while creating JWT token")
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	respondWithJSON(w, http.StatusOK, resp{
		Token: jwtToken,
	})
}
