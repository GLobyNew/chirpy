package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/GLobyNew/chirpy/internal/auth"
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
	structChirp := Chirp{
		ID:        foundChirp.ID,
		CreatedAt: foundChirp.CreatedAt,
		UpdatedAt: foundChirp.UpdatedAt,
		Body:      foundChirp.Body,
		UserID:    foundChirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, structChirp)

}

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {

	author_id := r.URL.Query().Get("author_id")
	sortQuery := r.URL.Query().Get("sort")
	if len(sortQuery) == 0 || (sortQuery != "desc" && sortQuery != "asc") {
		sortQuery = "asc"
	}

	var dbChirps []database.Chirp
	var err error

	// Check if author_id is not empty
	if len(author_id) == 0 {
		dbChirps, err = cfg.db.GetAllChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
	} else {
		authorID, err := uuid.Parse(author_id)
		if err != nil {
			log.Println("error while parsing author_id")
			respondWithError(w, http.StatusBadRequest, "invalid author ID")
			return
		}

		dbChirps, err = cfg.db.GetChirpsByUserID(r.Context(), authorID)
		if err != nil {
			log.Println("error while getting chirps from db")
			respondWithError(w, http.StatusInternalServerError, "error while getting chirps from db")
			return
		}
	}

	chirps, err := mapDatabaseChirpsToChirps(dbChirps)
	if err != nil {
		log.Println("error while mapping database chirps to chirps")
		respondWithError(w, http.StatusInternalServerError, "error while mapping database chirps to chirps")
		return
	}

	sort.Slice(chirps, func(i, j int) bool {
		if sortQuery == "asc" {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		}
		return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handleChirpCreation(w http.ResponseWriter, r *http.Request) {
	log.SetPrefix("handleChirpCreation: ")
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("can't get bearer token : %v", err)
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	userID, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)
	if err != nil {
		log.Printf("can't validate JWT: %v", err)
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	type parameters struct {
		Body string `json:"body"`
	}

	// Try decode request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
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
		UserID: userID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	structChirp := Chirp{
		ID:        createdChirp.ID,
		CreatedAt: createdChirp.CreatedAt,
		UpdatedAt: createdChirp.UpdatedAt,
		Body:      createdChirp.Body,
		UserID:    createdChirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, structChirp)

}

func (cfg *apiConfig) handleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
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

	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		log.Println("error while getting chirp from db")
		respondWithError(w, http.StatusNotFound, "chirp not found")
		return
	}
	if chirp.UserID != user.ID {
		log.Println("user id from db and token do not match")
		respondWithError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}
	err = cfg.db.DeleteChirpByID(r.Context(), chirpID)
	if err != nil {
		log.Println("error while deleting chirp from db")
		respondWithError(w, http.StatusInternalServerError, "error while deleting chirp from db")
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func mapDatabaseChirpsToChirps(dbChirps []database.Chirp) ([]Chirp, error) {
	chirps := make([]Chirp, len(dbChirps))
	for i, dbChirp := range dbChirps {
		chirps[i] = Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		}
	}
	return chirps, nil
}
