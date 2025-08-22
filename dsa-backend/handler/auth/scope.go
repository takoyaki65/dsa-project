package auth

import "fmt"

const (
	ScopeGrading = "grading"
	ScopeAdmin   = "admin"
)

const (
	RoleAdmin   = "admin"
	RoleManager = "manager"
	RoleStudent = "student"
)

func UserScopes() map[string][]string {
	return map[string][]string{
		RoleAdmin:   {ScopeGrading, ScopeAdmin},
		RoleManager: {ScopeGrading},
		RoleStudent: {},
	}
}

func UserRolesToID() map[string]int {
	return map[string]int{
		RoleAdmin:   1,
		RoleManager: 2,
		RoleStudent: 3,
	}
}

func RoleIDToUserRole() map[int]string {
	return map[int]string{
		1: RoleAdmin,
		2: RoleManager,
		3: RoleStudent,
	}
}

func GetUserScopes(userRole string) ([]string, error) {
	scopes, ok := UserScopes()[userRole]
	if !ok {
		return nil, fmt.Errorf("invalid user role: %s", userRole)
	}
	return scopes, nil
}

func GetRoleID(userRole string) (int, error) {
	roleID, ok := UserRolesToID()[userRole]
	if !ok {
		return 0, fmt.Errorf("invalid user role: %s", userRole)
	}
	return roleID, nil
}

func GetRoleName(roleID int) (string, error) {
	roleName, ok := RoleIDToUserRole()[roleID]
	if !ok {
		return "", fmt.Errorf("invalid role ID: %d", roleID)
	}
	return roleName, nil
}
