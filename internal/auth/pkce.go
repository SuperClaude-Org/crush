package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

type PKCEChallenge struct {
	Verifier  string
	Challenge string
	Method    string
}

func GeneratePKCEChallenge() (*PKCEChallenge, error) {
	verifier, err := generateCodeVerifier(64)
	if err != nil {
		return nil, fmt.Errorf("failed to generate code verifier: %w", err)
	}

	// Create SHA-256 challenge from verifier
	challenge := createCodeChallenge(verifier)

	return &PKCEChallenge{
		Verifier:  verifier,
		Challenge: challenge,
		Method:    "S256",
	}, nil
}

// generateCodeVerifier creates a cryptographically random code verifier
func generateCodeVerifier(length int) (string, error) {
	// Use URL-safe characters: [A-Z] / [a-z] / [0-9] / "-" / "." / "_" / "~"
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~"
	
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}

	return string(b), nil
}

// createCodeChallenge creates a SHA-256 base64url-encoded challenge from the verifier
func createCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(hash[:])
}