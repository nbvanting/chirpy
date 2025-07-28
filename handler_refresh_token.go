package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/nbvanting/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
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

	// Look up token in DB
	rt, err := cfg.db.GetRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", nil)
		return
	}
	if rt.RevokedAt.Valid || time.Now().After(rt.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token expired or revoked", nil)
		return
	}

	// Create new access token for the user
	accessToken, err := auth.MakeJWT(rt.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create access token", err)
		return
	}

	// Respond with the new access token
	respondWithJSON(w, http.StatusOK, map[string]string{"token": accessToken})

}
