package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/GLobyNew/chirpy/internal/auth"
	"github.com/GLobyNew/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Email         string    `json:"email"`
	Token         string    `json:"token"`
	Refresh_Token string    `json:"refresh_token"`
	Is_Chirpy_Red bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handleUser(w http.ResponseWriter, r *http.Request) {
	type inParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Try decode request
	decoder := json.NewDecoder(r.Body)
	params := inParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error while decoding request")
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Println("error while hashing password")
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPass,
	})
	if err != nil {
		log.Println("error while creating user in db")
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	userStruct := User{
		ID:            user.ID,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		Email:         user.Email,
		Is_Chirpy_Red: user.ChirpyRed,
	}

	respondWithJSON(w, http.StatusCreated, userStruct)

}

func (cfg *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	authToken, err := auth.GetBearerToken(r.Header)

	if err != nil {
		log.Println("bad auth token")
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return

	}

	userUUID, err := auth.ValidateJWT(authToken, cfg.jwtSecret)
	if err != nil {
		log.Println("fail while validating JWT")
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	user, err := cfg.db.GetUserByUUID(r.Context(), userUUID)
	if err != nil {
		log.Println("error while getting user from db")
		respondWithError(w, http.StatusInternalServerError, "error while getting user from db")
		return
	}
	if user.ID != userUUID {
		log.Println("user id from db and token do not match")
		respondWithError(w, http.StatusInternalServerError, "user id from db and token do not match")
		return
	}

	type inParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Try decode request
	decoder := json.NewDecoder(r.Body)
	params := inParams{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error while decoding request")
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Println("error while hashing password")
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	user, err = cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             user.ID,
		Email:          params.Email,
		HashedPassword: hashedPass,
	})
	if err != nil {
		log.Println("error while updating user in db")
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	userStruct := User{
		ID:            user.ID,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		Email:         user.Email,
		Is_Chirpy_Red: user.ChirpyRed,
	}

	respondWithJSON(w, http.StatusOK, userStruct)

}
