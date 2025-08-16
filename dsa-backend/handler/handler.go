package handler

import (
	"context"
	"crypto/rand"
	"dsa-backend/store"
	"encoding/hex"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

type Handler struct {
	userStore    store.UserStore
	jwtSecret    string
	shutdownChan chan struct{}
}

func NewHandler(db *bun.DB) *Handler {
	return &Handler{
		userStore:    *store.NewUserStore(db),
		jwtSecret:    generateSecretKey(),
		shutdownChan: make(chan struct{}),
	}
}

func (h *Handler) GetShutdownChan() chan struct{} {
	return h.shutdownChan
}

func (h *Handler) RegisterRoutes(r *echo.Echo) {
	// Check if the admin user exists, if not, only expose the admin creation endpoint
	ctx := context.Background()
	admins, err := h.userStore.GetUserListByUserRole(&ctx, "admin")
	if err != nil {
		panic(fmt.Sprintf("failed to check admin user: %v", err))
	}

	if len(*admins) == 0 {
		r.POST("/admin/create", h.CreateAdminUser)
		return
	}

	// If admin user exists, register all routes
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
