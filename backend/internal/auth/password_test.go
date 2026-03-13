package auth

import "testing"

func TestPasswordHashAndCheck(t *testing.T) {
	hash, err := HashPassword("Admin123!")
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}
	if hash == "" {
		t.Fatalf("expected non-empty hash")
	}
	if !CheckPassword(hash, "Admin123!") {
		t.Fatalf("expected password to match")
	}
	if CheckPassword(hash, "wrong") {
		t.Fatalf("expected password mismatch")
	}
}
