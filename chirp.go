package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/GLobyNew/chirpy/internal/database"
	"github.com/google/uuid"
)

const (
	maxChirpyLen        = 140
	errorGeneric        = "Something went wrong"
	errorChirpIsTooLong = "Chirp is too long"
	errorBodyIsEmpty    = "Body is empty"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func unProfaneChirp(s string) string {
	var forbiddenWords = [...]string{"kerfuffle", "sharbert", "fornax"}
	splittedStr := strings.Split(s, " ")
	for i, word := range splittedStr {
		for _, prWord := range forbiddenWords {
			if strings.ToLower(word) == prWord {
				replaceWord := "****"
				splittedStr[i] = replaceWord
			}
		}
	}
	return strings.Join(splittedStr, " ")
}

func (cfg *apiConfig) handleChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	// Try decode request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error while decoding request")
		return
	}

	chirp := unProfaneChirp(params.Body)

	// Check if not empty
	if len(chirp) == 0 {
		respondWithError(w, http.StatusBadRequest, "chirp is empty")
		return
	}

	// Check max against max length
	if len(chirp) > maxChirpyLen {
		respondWithError(w, http.StatusBadRequest, "chirp lenght is more than 140 symbols")
		return
	}

	createdChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   chirp,
		UserID: params.UserID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        createdChirp.ID,
		CreatedAt: createdChirp.CreatedAt,
		UpdatedAt: createdChirp.UpdatedAt,
		Body:      createdChirp.Body,
		UserID:    createdChirp.UserID,
	})

}
