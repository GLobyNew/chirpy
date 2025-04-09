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
	ID         uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Email      string    `json:"email"`
	hashedPass string
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
	userStruct, err := mapDatabaseUserToUserStruct(user)
	if err != nil {
		log.Println("error while converting user db to user struct in handleUser func")
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	respondWithJSON(w, http.StatusCreated, userStruct)

}

func mapDatabaseUserToUserStruct(dbUser database.User) (User, error) {
	// Marshal the database chirps into JSON
	data, err := json.Marshal(dbUser)
	if err != nil {
		return User{}, err
	}

	// Unmarshal the JSON into the main Chirp struct
	var user User
	err = json.Unmarshal(data, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
