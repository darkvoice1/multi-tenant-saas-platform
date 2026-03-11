package auth

const (
	RoleAdmin   = "admin"
	RoleManager = "manager"
	RoleMember  = "member"
	RoleGuest   = "guest"
)

var roleRank = map[string]int{
	RoleGuest:   1,
	RoleMember:  2,
	RoleManager: 3,
	RoleAdmin:   4,
}

func IsRoleAtLeast(role, required string) bool {
	return roleRank[role] >= roleRank[required]
}

func IsValidRole(role string) bool {
	_, ok := roleRank[role]
	return ok
}
