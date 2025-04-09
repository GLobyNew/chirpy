package main

import (
	"encoding/json"
	"log"
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
	errorChirpIsEmpty   = "Chirp is empty"
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

func (cfg *apiConfig) handleGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		log.Println("error while parsing UUID in path value in GetChirp")
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	foundChirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		log.Println("No chirp was found")
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}
	structChirp, err := mapDatabaseChirpToChirps(foundChirp)
	if err != nil {
		log.Println("error while marshalling/unmarshalling db chirp to go struct chirp")
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	respondWithJSON(w, http.StatusOK, structChirp)

}

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	chirps, err := mapDatabaseChirpsToChirps(dbChirps)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handleChirpCreation(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	// Try decode request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	chirp := unProfaneChirp(params.Body)

	// Check if not empty
	if len(chirp) == 0 {
		respondWithError(w, http.StatusBadRequest, errorChirpIsEmpty)
	}

	// Check max against max length
	if len(chirp) > maxChirpyLen {
		respondWithError(w, http.StatusBadRequest, errorChirpIsTooLong)
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

	structChirp, err := mapDatabaseChirpToChirps(createdChirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	respondWithJSON(w, http.StatusCreated, structChirp)

}

func mapDatabaseChirpsToChirps(dbChirps []database.Chirp) ([]Chirp, error) {
	// Marshal the database chirps into JSON
	data, err := json.Marshal(dbChirps)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON into the main Chirp struct
	var chirps []Chirp
	err = json.Unmarshal(data, &chirps)
	if err != nil {
		return nil, err
	}

	return chirps, nil
}

func mapDatabaseChirpToChirps(dbChirp database.Chirp) (Chirp, error) {
	// Marshal the database chirps into JSON
	data, err := json.Marshal(dbChirp)
	if err != nil {
		return Chirp{}, err
	}

	// Unmarshal the JSON into the main Chirp struct
	var chirp Chirp
	err = json.Unmarshal(data, &chirp)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}
