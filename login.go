package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/GLobyNew/chirpy/internal/auth"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
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

	foundUserStruct, err := mapDatabaseUserToUserStruct(foundUser)
	if err != nil {
		log.Println("error while converting user db to user struct in handleUser func")
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	respondWithJSON(w, http.StatusOK, foundUserStruct)

}
