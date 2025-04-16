package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/GLobyNew/chirpy/internal/auth"
	"github.com/GLobyNew/chirpy/internal/database"
)

const (
	DefaultExpiresIn = time.Hour
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	log.SetPrefix("handleLogin: ")
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	// Try decode request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	foundUser, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Println("no user was found in handleLogin")
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	err = auth.CheckPasswordHash(foundUser.HashedPassword, params.Password)
	if err != nil {
		log.Println("password is incorrect")
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	// Complete the User struct initialization
	userStruct := User{
		ID:            foundUser.ID,
		CreatedAt:     foundUser.CreatedAt,
		UpdatedAt:     foundUser.UpdatedAt,
		Email:         foundUser.Email,
		Is_Chirpy_Red: foundUser.ChirpyRed,
	}

	// Generate JWT token
	jwtToken, err := auth.MakeJWT(foundUser.ID, cfg.jwtSecret, DefaultExpiresIn)
	if err != nil {
		log.Println("error while creating JWT token")
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	userStruct.Token = jwtToken

	// Generate refresh token
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Println("error while creating refresh token")
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	// Store refresh token in the database
	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    userStruct.ID,
		ExpiresAt: time.Now().Add(1440 * time.Hour),
	})
	if err != nil {
		log.Println("error while adding refresh token to db")
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	userStruct.Refresh_Token = refreshToken

	// Respond with the user data
	respondWithJSON(w, http.StatusOK, userStruct)

}
