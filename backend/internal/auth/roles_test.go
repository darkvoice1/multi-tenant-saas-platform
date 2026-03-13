package auth

import "testing"

func TestRoles(t *testing.T) {
	if !IsValidRole(RoleAdmin) || !IsValidRole(RoleMember) {
		t.Fatalf("expected valid roles")
	}
	if IsValidRole("unknown") {
		t.Fatalf("expected invalid role")
	}
	if !IsRoleAtLeast(RoleAdmin, RoleMember) {
		t.Fatalf("expected admin to be >= member")
	}
	if IsRoleAtLeast(RoleGuest, RoleAdmin) {
		t.Fatalf("expected guest to be < admin")
	}
}
