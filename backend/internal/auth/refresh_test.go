package auth

import "testing"

func TestRefreshTokenGeneration(t *testing.T) {
	token, hash, err := GenerateRefreshToken()
	if err != nil {
		t.Fatalf("GenerateRefreshToken error: %v", err)
	}
	if token == "" || hash == "" {
		t.Fatalf("expected token and hash")
	}
	if HashRefreshToken(token) != hash {
		t.Fatalf("hash mismatch")
	}
}
