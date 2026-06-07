package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Khaz713/chirpy/internal/database"
	"github.com/google/uuid"
)

type parametersWebhook struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := &parametersWebhook{}
	err := decoder.Decode(params)
	if err != nil {
		log.Printf("Error decoding body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if params.Event != "user.upgraded" {
		log.Printf("Event not supported: %v", params.Event)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		log.Printf("Error parsing user id: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = cfg.db.UserSetRed(r.Context(), database.UserSetRedParams{
		ID:          userID,
		IsChirpyRed: true,
	})
	if err != nil {
		log.Printf("Error setting user as red: %v", err)
		w.WriteHeader(http.StatusNotFound)
	}
	w.WriteHeader(http.StatusNoContent)
}
