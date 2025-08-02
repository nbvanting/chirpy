package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nbvanting/chirpy/internal/database"
)

// handlerListAllChirps handles GET requests to list all chirps.
func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authorIDStr := r.URL.Query().Get("author_id")
	var chirps []database.Chirp
	var err error
	if authorIDStr != "" {
		authorUUID, errParse := uuid.Parse(authorIDStr)
		if errParse != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid author_id"})
			return
		}
		chirps, err = cfg.db.GetChirpsByAuthor(ctx, authorUUID)
	} else {
		chirps, err = cfg.db.GetChirps(ctx)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "could not fetch chirps"})
		return
	}

	// Map to response structure used by POST /api/chirps
	type chirpResponse struct {
		ID        string `json:"id"`
		Body      string `json:"body"`
		AuthorID  string `json:"author_id"`
		CreatedAt string `json:"created_at"`
	}

	resp := make([]chirpResponse, 0, len(chirps))
	for _, c := range chirps {
		resp = append(resp, chirpResponse{
			ID:        c.ID.String(),
			Body:      c.Body,
			AuthorID:  c.UserID.String(),
			CreatedAt: c.CreatedAt.Format(time.RFC3339),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// handlerChirpRetrieveByID handles GET requests to fetch a single chirp by its ID.
func (cfg *apiConfig) handlerChirpRetrieveByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	chirpIDStr := r.PathValue("chirpID")

	chirpUUID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid chirp id"})
		return
	}

	chirp, err := cfg.db.GetChirp(ctx, chirpUUID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "chirp not found"})
		return
	}

	resp := map[string]interface{}{
		"id":         chirp.ID.String(),
		"created_at": chirp.CreatedAt.Format(time.RFC3339),
		"updated_at": chirp.UpdatedAt.Format(time.RFC3339),
		"body":       chirp.Body,
		"user_id":    chirp.UserID.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
