package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/nbvanting/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds *int   `json:"expires_in_seconds,omitempty"`
	}
	type response struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Email     string `json:"email"`
		Token     string `json:"token"`
		ExpiresAt string `json:"expires_at"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	dbUser, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	err = auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	// Determine expiration
	const maxExpiration = 3600 // 1 hour in seconds
	expirationSeconds := maxExpiration
	if params.ExpiresInSeconds != nil {
		if *params.ExpiresInSeconds < maxExpiration {
			expirationSeconds = *params.ExpiresInSeconds
		}
	}

	// Generate token using MakeJWT
	tokenSecret := cfg.jwtSecret
	token, err := auth.MakeJWT(dbUser.ID, tokenSecret, time.Duration(expirationSeconds)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not generate token", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		ID:        dbUser.ID.String(),
		CreatedAt: dbUser.CreatedAt.Format(time.RFC3339),
		UpdatedAt: dbUser.UpdatedAt.Format(time.RFC3339),
		Email:     dbUser.Email,
		Token:     token,
	})
}
