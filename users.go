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
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

type parametersUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type parametersLogin struct {
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
	login := parametersLogin{}
	err := decoder.Decode(&login)
	if err != nil {
		log.Printf("Error decoding login parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error decoding login parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), login.Email)
	if err != nil {
		log.Printf("Error getting user by email: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error getting user by email", err)
		return
	}
	passCorrect, err := auth.CheckPasswordHash(login.Password, user.HashedPassword)
	if err != nil {
		log.Printf("Error checking password hash: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error checking password hash", err)
		return
	}

	if !passCorrect {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Second*time.Duration(3600))
	if err != nil {
		log.Printf("Error generating token: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error generating token", err)
		return
	}

	refreshToken := auth.MakeRefreshToken()

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
		UserID:    user.ID,
	})
	if err != nil {
		log.Printf("Error creating refresh token: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error creating refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
	})

}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	newUser := parametersUser{}
	err := decoder.Decode(&newUser)
	if err != nil {
		log.Printf("Error decoding user parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error decoding user parameters", err)
		return
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %s", err)
		respondWithError(w, http.StatusUnauthorized, "Error getting bearer token", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error validating bearer token: %s", err)
		respondWithError(w, http.StatusUnauthorized, "Error validating bearer token", err)
		return
	}
	hashedPass, err := auth.HashPassword(newUser.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}
	updatedUser, err := cfg.db.UpdateUserByID(r.Context(), database.UpdateUserByIDParams{
		ID:             userID,
		Email:          newUser.Email,
		HashedPassword: hashedPass,
	})
	if err != nil {
		log.Printf("Error updating user: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error updating user", err)
		return
	}
	respondWithJSON(w, http.StatusOK, User{
		ID:        updatedUser.ID,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email:     updatedUser.Email,
	})
}
