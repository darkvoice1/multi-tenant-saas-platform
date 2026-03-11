package auth

type Permission string

const (
	PermTenantRead   Permission = "tenant:read"
	PermTenantWrite  Permission = "tenant:write"
	PermUserRead     Permission = "user:read"
	PermUserWrite    Permission = "user:write"
	PermProjectRead  Permission = "project:read"
	PermProjectWrite Permission = "project:write"
	PermTaskRead     Permission = "task:read"
	PermTaskWrite    Permission = "task:write"
	PermAuditRead    Permission = "audit:read"
	PermAuditWrite   Permission = "audit:write"
	PermAdminPing    Permission = "admin:ping"
)

var PermissionMatrix = map[Permission][]string{
	PermTenantRead:   {RoleAdmin, RoleManager},
	PermTenantWrite:  {RoleAdmin},
	PermUserRead:     {RoleAdmin, RoleManager},
	PermUserWrite:    {RoleAdmin},
	PermProjectRead:  {RoleAdmin, RoleManager, RoleMember, RoleGuest},
	PermProjectWrite: {RoleAdmin, RoleManager},
	PermTaskRead:     {RoleAdmin, RoleManager, RoleMember, RoleGuest},
	PermTaskWrite:    {RoleAdmin, RoleManager, RoleMember},
	PermAuditRead:    {RoleAdmin},
	PermAuditWrite:   {RoleAdmin},
	PermAdminPing:    {RoleAdmin},
}

func IsAllowed(role string, perm Permission) bool {
	allowed, ok := PermissionMatrix[perm]
	if !ok {
		return false
	}
	for _, r := range allowed {
		if r == role {
			return true
		}
	}
	return false
}
