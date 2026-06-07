package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Khaz713/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error getting bearer token", err)
	}
	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		log.Printf("Error getting user from refresh token: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Token does not exist/expired", err)
		return
	}
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Second*time.Duration(3600))
	if err != nil {
		log.Printf("Error generating token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error generating token", err)
		return
	}
	type body struct {
		Token string `json:"token"`
	}
	respondWithJSON(w, http.StatusOK, body{Token: token})
}

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error getting bearer token", err)
	}
	err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		log.Printf("Error revoking refresh token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error revoking refresh token", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
