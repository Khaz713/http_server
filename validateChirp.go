package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type parameters struct {
	Body string `json:"body"`
}

type returnError struct {
	Error string `json:"error"`
}

type returnCleaned struct {
	CleanedBody string `json:"cleaned_body"`
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", err)
		return
	}
	message := profanitiesCheck(params.Body)

	respondWithJSON(w, http.StatusOK, returnCleaned{CleanedBody: message})

}

func profanitiesCheck(message string) string {
	profanities := map[string]struct{}{
		"fornax":    {},
		"kerfuffle": {},
		"sharbert":  {},
	}
	words := strings.Split(message, " ")
	for i, word := range words {
		lowerWord := strings.ToLower(word)
		if _, ok := profanities[lowerWord]; ok {
			words[i] = "****"
		}
	}
	cleanedMessage := strings.Join(words, " ")

	return cleanedMessage
}
