package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestHashPasswordAndCheckPasswordHash(t *testing.T) {
	password := "mySecretPassword123!"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if hash == "" {
		t.Error("HashPassword returned empty hash")
	}

	// Check that the password matches the hash
	err = CheckPasswordHash(password, hash)
	if err != nil {
		t.Errorf("CheckPasswordHash failed for correct password: %v", err)
	}

	// Check that a wrong password does not match
	wrongPassword := "notMyPassword"
	err = CheckPasswordHash(wrongPassword, hash)
	if err == nil {
		t.Error("CheckPasswordHash did not fail for incorrect password")
	}
}

func TestMakeJWTAndValidateJWT_Success(t *testing.T) {
	userID := uuid.New()
	secret := "testsecret"
	expires := time.Minute
	token, err := MakeJWT(userID, secret, expires)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	parsedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}
	if parsedID != userID {
		t.Errorf("Expected userID %v, got %v", userID, parsedID)
	}
}

func TestValidateJWT_InvalidSignature(t *testing.T) {
	userID := uuid.New()
	secret := "testsecret"
	wrongSecret := "wrongsecret"
	token, err := MakeJWT(userID, secret, time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}
	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Error("Expected error for invalid signature, got nil")
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "testsecret"
	token, err := MakeJWT(userID, secret, -time.Minute) // already expired
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}
	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Error("Expected error for expired token, got nil")
	}
}

func TestValidateJWT_MalformedToken(t *testing.T) {
	secret := "testsecret"
	malformed := "not.a.jwt.token"
	_, err := ValidateJWT(malformed, secret)
	if err == nil {
		t.Error("Expected error for malformed token, got nil")
	}
}

func TestValidateJWT_InvalidUUIDSubject(t *testing.T) {
	secret := "testsecret"
	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute)),
		Subject:   "not-a-uuid",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}
	_, err = ValidateJWT(tokenStr, secret)
	if err == nil {
		t.Error("Expected error for invalid UUID subject, got nil")
	}
}
