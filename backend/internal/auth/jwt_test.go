package auth

import (
	"testing"
	"time"
)

func TestCreateAndParseAccessToken(t *testing.T) {
	secret := "test_secret"
	userID := "user-1"
	tenantID := "tenant-1"
	role := RoleAdmin

	token, err := CreateAccessToken(secret, time.Minute, userID, tenantID, role)
	if err != nil {
		t.Fatalf("CreateAccessToken error: %v", err)
	}

	claims, err := ParseToken(secret, token)
	if err != nil {
		t.Fatalf("ParseToken error: %v", err)
	}
	if claims.UserID != userID || claims.TenantID != tenantID || claims.Role != role {
		t.Fatalf("claims mismatch: %+v", claims)
	}
}

func TestParseTokenFailures(t *testing.T) {
	secret := "test_secret"

	if _, err := ParseToken(secret, "not-a-token"); err == nil {
		t.Fatalf("expected error for invalid token")
	}

	valid, err := CreateAccessToken(secret, time.Minute, "u", "t", RoleMember)
	if err != nil {
		t.Fatalf("CreateAccessToken error: %v", err)
	}

	if _, err := ParseToken("wrong_secret", valid); err == nil {
		t.Fatalf("expected error for wrong secret")
	}
}
