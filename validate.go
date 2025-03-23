package main

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	maxChirpyLen        = 140
	errorGeneric        = "Something went wrong"
	errorChirpIsTooLong = "Chirp is too long"
	errorBodyIsEmpty    = "Body is empty"
)

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnValid struct {
		Valid bool `json:"valid"`
	}
	type returnError struct {
		ErrorStr string `json:"error"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		data, err := json.Marshal(returnError{
			ErrorStr: errorGeneric,
		})
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(data)
		return
	}
	if len(params.Body) == 0 {
		data, err := json.Marshal(returnError{
			ErrorStr: errorBodyIsEmpty,
		})
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(data)
		return
	}
	if len(params.Body) > maxChirpyLen {
		data, err := json.Marshal(returnError{
			ErrorStr: errorChirpIsTooLong,
		})
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(data)
		return
	}

	data, err := json.Marshal(returnValid{
		Valid: true,
	})
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
