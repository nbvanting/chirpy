package auth

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

// MakeRefreshToken generates a random 256-bit (32-byte) hex-encoded string.
func MakeRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
