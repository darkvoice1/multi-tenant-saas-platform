package auth

import "testing"

func TestIsAllowed(t *testing.T) {
	cases := []struct {
		name    string
		role    string
		perm    Permission
		allowed bool
	}{
		{"admin tenant read", RoleAdmin, PermTenantRead, true},
		{"manager tenant read", RoleManager, PermTenantRead, true},
		{"member tenant read", RoleMember, PermTenantRead, false},
		{"guest tenant read", RoleGuest, PermTenantRead, false},
		{"admin tenant write", RoleAdmin, PermTenantWrite, true},
		{"manager tenant write", RoleManager, PermTenantWrite, false},
		{"admin user read", RoleAdmin, PermUserRead, true},
		{"manager user read", RoleManager, PermUserRead, true},
		{"member user read", RoleMember, PermUserRead, false},
		{"admin user write", RoleAdmin, PermUserWrite, true},
		{"manager user write", RoleManager, PermUserWrite, false},
		{"guest project read", RoleGuest, PermProjectRead, true},
		{"member project read", RoleMember, PermProjectRead, true},
		{"manager project read", RoleManager, PermProjectRead, true},
		{"guest project write", RoleGuest, PermProjectWrite, false},
		{"member project write", RoleMember, PermProjectWrite, false},
		{"manager project write", RoleManager, PermProjectWrite, true},
		{"member task read", RoleMember, PermTaskRead, true},
		{"guest task read", RoleGuest, PermTaskRead, true},
		{"guest task write", RoleGuest, PermTaskWrite, false},
		{"member task write", RoleMember, PermTaskWrite, true},
		{"manager task write", RoleManager, PermTaskWrite, true},
		{"admin audit read", RoleAdmin, PermAuditRead, true},
		{"manager audit read", RoleManager, PermAuditRead, false},
		{"admin audit write", RoleAdmin, PermAuditWrite, true},
		{"member audit write", RoleMember, PermAuditWrite, false},
		{"admin ping", RoleAdmin, PermAdminPing, true},
		{"guest ping", RoleGuest, PermAdminPing, false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if got := IsAllowed(tc.role, tc.perm); got != tc.allowed {
				t.Fatalf("expected %v, got %v", tc.allowed, got)
			}
		})
	}
}
