package auth

import (
	"testing"
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
