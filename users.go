package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	type Email struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	email := Email{}
	err := decoder.Decode(&email)
	if err != nil {
		log.Printf("Error decoding email: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error decoding email", err)
		return
	}
	user, err := cfg.db.CreateUser(r.Context(), email.Email)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error creating user", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
