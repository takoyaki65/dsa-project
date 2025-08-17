package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type JwtCustomClaims struct {
	UserID string   `json:"userid"`
	Scopes []string `json:"scopes"`
	jwt.RegisteredClaims
}

func IssueNewToken(userid string, scopes []string, secret string, issuedAt time.Time) (string, error) {
	newClaim := newJwtCustomClaims(userid, scopes, issuedAt)
	newToken := createToken(newClaim)
	newTokenStr, err := newToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return newTokenStr, nil
}

func createToken(claim *JwtCustomClaims) *jwt.Token {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
}

func newJwtCustomClaims(userid string, scopes []string, issuedAt time.Time) *JwtCustomClaims {
	return &JwtCustomClaims{
		UserID: userid,
		Scopes: scopes,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(issuedAt.Add(time.Hour * 12)),
		},
	}
}
