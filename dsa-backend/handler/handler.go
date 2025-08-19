package handler

import (
	"context"
	"crypto/rand"
	"dsa-backend/storage"
	"dsa-backend/utils"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

type Handler struct {
	userStore storage.UserStore
	jwtSecret string
}

func NewHandler(db *bun.DB) *Handler {
	return &Handler{
		userStore: *storage.NewUserStore(db),
		jwtSecret: generateSecretKey(),
	}
}

func (h *Handler) RegisterRoutes(r *echo.Group) {
	r.POST("/login", h.Login)

	currentUserGroup := r.Group("/currentUser")
	currentUserGroup.Use(echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(utils.JwtCustomClaims)
		},
		SigningKey: []byte(h.jwtSecret),
	}))
	currentUserGroup.Use(h.CheckValidityOfJWTMiddleware())
	currentUserGroup.Use(RequiredScopesMiddleware("me"))
	currentUserGroup.GET("/me", h.GetCurrentUser)

}

func generateSecretKey() string {
	// generate random 32-byte key
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic(fmt.Sprintf("failed to generate secret key: %v", err))
	}
	// convert to hex
	return hex.EncodeToString(key)
}

func (h *Handler) CheckValidityOfJWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, ok := c.Get("user").(*jwt.Token)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			claims, ok := token.Claims.(*utils.JwtCustomClaims)

			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
			}

			// Check token expiration
			if claims.ExpiresAt.Time.Before(time.Now()) {
				return echo.NewHTTPError(http.StatusUnauthorized, "token expired")
			}

			ctx := context.Background()
			// Check login history existence
			loginHistory, err := h.userStore.GetLoginHistory(&ctx, claims.UserID, claims.IssuedAt.Time)

			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to get login history")
			}

			if loginHistory == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "login history not found")
			}

			return next(c)
		}
	}
}
