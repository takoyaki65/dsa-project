package handler

import "fmt"

func UserScopes() map[string][]string {
	return map[string][]string{
		"admin":   {"me", "grading", "admin"},
		"manager": {"me", "grading"},
		"student": {"me"},
	}
}

func UserRolesToID() map[string]int {
	return map[string]int{
		"admin":   1,
		"manager": 2,
		"student": 3,
	}
}

func RoleIDToUserRole() map[int]string {
	return map[int]string{
		1: "admin",
		2: "manager",
		3: "student",
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
