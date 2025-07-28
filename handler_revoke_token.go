package main

import (
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	// Get Authorization header
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid Authorization header", nil)
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	token = strings.TrimSpace(token)
	if token == "" {
		respondWithError(w, http.StatusUnauthorized, "Missing refresh token", nil)
		return
	}

	// Revoke the token in the database
	err := cfg.db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", nil)
		return
	}

	// Respond with 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
