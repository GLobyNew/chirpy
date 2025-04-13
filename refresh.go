package main

import (
	"log"
	"net/http"
	"time"

	"github.com/GLobyNew/chirpy/internal/auth"
	"github.com/GLobyNew/chirpy/internal/database"
)

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, r *http.Request) {

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

	if tokenInDB.RevokedAt.Valid && tokenInDB.RevokedAt.Time.Before((time.Now())) {
		log.Println("token is revoked")
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	newToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Println("error while creating new token")
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	createdToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     newToken,
		UserID:    tokenInDB.UserID,
		ExpiresAt: time.Now().Add(1440 * time.Hour),
	})
	if err != nil {
		log.Println("error while adding refresh token to db")
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	respondWithJSON(w, http.StatusOK, resp{
		Token: createdToken.Token,
	})
}
