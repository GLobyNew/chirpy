package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func (cfg *apiConfig) handleUser(w http.ResponseWriter, r *http.Request) {
	type inParams struct {
		Email string `json:"email"`
	}

	type response struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
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

	payload := response{
		ID:        user.ID.String(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, http.StatusOK, payload)

}
