package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/nbvanting/chirpy/internal/auth"
	"github.com/nbvanting/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		ID           string `json:"id"`
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
		Email        string `json:"email"`
		Token        string `json:"token"`
		ExpiresAt    string `json:"expires_at"`
		RefreshToken string `json:"refresh_token"`
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

	// Access token: 1 hour expiration
	accessTokenDuration := time.Hour
	tokenSecret := cfg.jwtSecret
	token, err := auth.MakeJWT(dbUser.ID, tokenSecret, accessTokenDuration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not generate token", nil)
		return
	}

	// Refresh token: 60 days expiration
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not generate refresh token", nil)
		return
	}

	// Store refresh token in database
	refreshTokenExpiration := time.Now().Add(60 * 24 * time.Hour) // 60 days
	err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    dbUser.ID,
		ExpiresAt: refreshTokenExpiration,
		// revoked_at will be null by default
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not store refresh token", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		ID:           dbUser.ID.String(),
		CreatedAt:    dbUser.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    dbUser.UpdatedAt.Format(time.RFC3339),
		Email:        dbUser.Email,
		Token:        token,
		ExpiresAt:    time.Now().Add(accessTokenDuration).Format(time.RFC3339),
		RefreshToken: refreshToken,
	})
}
