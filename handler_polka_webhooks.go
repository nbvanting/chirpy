package main

import (
	"encoding/json"
	"net/http"

	"database/sql"

	"github.com/google/uuid"
	"github.com/nbvanting/chirpy/internal/auth"
)

type WebhookPayload struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload WebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	apikey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if cfg.polkaKey != apikey {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if payload.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	uid, err := uuid.Parse(payload.Data.UserID)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	_, err = cfg.db.UpgradeUserToChirpyRed(r.Context(), uid)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
