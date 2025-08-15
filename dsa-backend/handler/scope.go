package handler

import "fmt"

func UserScopes() map[string][]string {
	return map[string][]string{
		"admin":   {"me", "grading", "admin"},
		"manager": {"me", "grading"},
		"student": {"me"},
	}
}

func GetUserScopes(userRole string) ([]string, error) {
	scopes, ok := UserScopes()[userRole]
	if !ok {
		return nil, fmt.Errorf("invalid user role: %s", userRole)
	}
	return scopes, nil
}
