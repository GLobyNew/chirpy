package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handleUser(w http.ResponseWriter, r *http.Request) {
	type inParams struct {
		Email string `json:"email"`
	}

	// Try decode request
	decoder := json.NewDecoder(r.Body)
	params := inParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error while decoding request")
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email)

	payload := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, http.StatusCreated, payload)

}
