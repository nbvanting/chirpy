package main

import (
	"encoding/json"
	"net/http"
	"time"
)

// handlerListAllChirps handles GET requests to list all chirps.
func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	chirps, err := cfg.db.GetChirps(ctx)
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
