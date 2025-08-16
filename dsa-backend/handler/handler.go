package handler

import (
	"context"
	"crypto/rand"
	"dsa-backend/store"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

type Handler struct {
	userStore    store.UserStore
	jwtSecret    string
	shutdownChan chan struct{}
	initialized  bool
}

func NewHandler(db *bun.DB) *Handler {
	return &Handler{
		userStore:    *store.NewUserStore(db),
		jwtSecret:    generateSecretKey(),
		shutdownChan: make(chan struct{}),
		initialized:  false,
	}
}

func (h *Handler) GetShutdownChan() chan struct{} {
	return h.shutdownChan
}

func (h *Handler) RegisterRoutes(r *echo.Group) {
	// Check if the admin user exists, if not, only expose the admin creation endpoint
	ctx := context.Background()
	admins, err := h.userStore.GetUserListByUserRole(&ctx, "admin")
	if err != nil {
		panic(fmt.Sprintf("failed to check admin user: %v", err))
	}

	if len(*admins) == 0 {
		r.POST("/admin/create", h.CreateAdminUser)
		h.initialized = false
		r.GET("/initialized", h.IsInitialized)
		return
	}

	h.initialized = true
	r.GET("/initialized", h.IsInitialized)

	// If admin user exists, register all routes
	r.POST("/login", h.Login)
}

// Initialized godoc
// @Summary Check if the application is initialized.
// @Descrition Returns whether the application has been initialized with an admin user.
// @Tags Initialization
// @Produce json
// @Success 200 {object} initCheckResponse "Returns a JSON object with the key 'initialized' set to true or false."
// @Router /initialized [get]
func (h *Handler) IsInitialized(c echo.Context) error {
	return c.JSON(http.StatusOK, newInitCheckResponse(h.initialized))
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
