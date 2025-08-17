package handler

import (
	"crypto/rand"
	"dsa-backend/store"
	"encoding/hex"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

type Handler struct {
	userStore store.UserStore
	jwtSecret string
}

func NewHandler(db *bun.DB) *Handler {
	return &Handler{
		userStore: *store.NewUserStore(db),
		jwtSecret: generateSecretKey(),
	}
}

func (h *Handler) RegisterRoutes(r *echo.Group) {
	r.POST("/login", h.Login)
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
