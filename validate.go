package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

const (
	maxChirpyLen        = 140
	errorGeneric        = "Something went wrong"
	errorChirpIsTooLong = "Chirp is too long"
	errorBodyIsEmpty    = "Body is empty"
)

func unProfaneChirp(s string) string {
	var forbiddenWords = [...]string{"kerfuffle", "sharbert", "fornax"}
	splittedStr := strings.Split(s, " ")
	for i, word := range splittedStr {
		for _, prWord := range forbiddenWords {
			if strings.ToLower(word) == prWord {
				replaceWord := strings.Repeat("*", len([]rune(splittedStr[i])))
				splittedStr[i] = replaceWord
			}
		}
	}
	return strings.Join(splittedStr, " ")
}

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}
	type clBody struct {
		ClearedBody string `json:"cleaned_body"`
	}

	// Try decode request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error while decoding request")
		return
	}

	chirp := params.Body

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

	// Check against forbidden words
	payload := clBody{
		ClearedBody: unProfaneChirp(chirp),
	}
	respondWithJSON(w, http.StatusOK, payload)

}
