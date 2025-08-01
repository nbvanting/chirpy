package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/nbvanting/chirpy/internal/auth"
	"github.com/nbvanting/chirpy/internal/database"
)

func (cfg *apiConfig) handlerChirpDelete(w http.ResponseWriter, r *http.Request) {
	// ENsure method is DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and validate access token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil || token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil || userID == uuid.Nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get chirpID from URL path, assuming route is /api/chirps/{chirpID}
	chirpIDStr := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		http.Error(w, "Chirp not found", http.StatusNotFound)
		return
	}

	// Check if user owns the chirp
	ownsChirp, err := cfg.db.CheckChirpOwnership(r.Context(), database.CheckChirpOwnershipParams{
		UserID:  userID,
		ID: chirpUUID,
	})
	if ownsChirp == uuid.Nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Delete the chirp
	if err := cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID:     chirpUUID,
		UserID: userID,
	}); err != nil {
		http.Error(w, "Error deleting chirp", http.StatusInternalServerError)
		return
	}

	// Respond with no content
	w.WriteHeader(http.StatusNoContent)
}
