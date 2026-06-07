package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Khaz713/chirpy/internal/auth"
	"github.com/Khaz713/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type parametersUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	email := parametersUser{}
	err := decoder.Decode(&email)
	if err != nil {
		log.Printf("Error decoding user parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error decoding user parameters", err)
		return
	}
	hashPass, err := auth.HashPassword(email.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          email.Email,
		HashedPassword: hashPass,
	})
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

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	email := parametersUser{}
	err := decoder.Decode(&email)
	if err != nil {
		log.Printf("Error decoding user parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error decoding user parameters", err)
		return
	}
	user, err := cfg.db.GetUserByEmail(r.Context(), email.Email)
	if err != nil {
		log.Printf("Error getting user by email: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error getting user by email", err)
		return
	}
	passCorrect, err := auth.CheckPasswordHash(email.Password, user.HashedPassword)
	if err != nil {
		log.Printf("Error checking password hash: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error checking password hash", err)
		return
	}

	if !passCorrect {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})

}
