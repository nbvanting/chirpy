package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/nbvanting/chirpy/internal/auth"
	"github.com/nbvanting/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {

	type response struct {
		User
	}

	// Ensure method is PUT
	if r.Method != http.MethodPut {
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

	// Get new email and password from JSON body
	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var reqBody requestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	newEmail := reqBody.Email
	newPassword := reqBody.Password
	if newEmail == "" || newPassword == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}
	// Hash the new password
	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Update user in DB
	params := database.UpdateUserParams{
		ID:             userID,
		Email:          newEmail,
		HashedPassword: hashedPassword,
	}
	user, err := cfg.db.UpdateUser(r.Context(), params)
	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	// Respond with updated user (without password)
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
