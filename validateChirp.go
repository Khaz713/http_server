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
	profanities := []string{"kerfuffle", "sharbert", "fornax"}
	cleanedMessage := message
	for _, profanity := range profanities {
		splitMessage := strings.Split(cleanedMessage, " ")
		for i, word := range splitMessage {
			if strings.ToLower(word) == strings.ToLower(profanity) {
				splitMessage[i] = "****"
			}
		}
		cleanedMessage = strings.Join(splitMessage, " ")
	}

	return cleanedMessage
}
