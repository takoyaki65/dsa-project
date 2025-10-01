package admin

import (
	"context"
	"dsa-backend/handler/response"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/takoyaki65/dsa-project/database/model/userrole"
)

type archiveUserRequest struct {
	UserID string `param:"user_id" validate:"required"`
}

func (h *Handler) ArchiveUser(c echo.Context) error {
	ctx := context.Background()

	var req archiveUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewError("Invalid request body: "+err.Error()))
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewError("Validation failed: "+err.Error()))
	}

	// Check if user exists
	user, err := h.userStore.GetUserByUserID(ctx, req.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewError("Failed to check user existence: "+err.Error()))
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, response.NewError("User not found"))
	}

	// Validation: Prevent archiving the last admin user
	if user.RoleID == userrole.Admin {
		return c.JSON(http.StatusForbidden, response.NewError("Cannot archive an admin user"))
	}

	// Modify user's DisabledAt to current time
	if err := h.userStore.ModifyUserValidity(ctx, req.UserID, time.Now()); err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewError("Failed to archive user: "+err.Error()))
	}

	return c.JSON(http.StatusOK, response.NewSuccess("User archived successfully"))
}

func (h *Handler) ActivateUser(c echo.Context) error {
	ctx := context.Background()

	var req archiveUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewError("Invalid request body: "+err.Error()))
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewError("Validation failed: "+err.Error()))
	}

	// Check if user exists
	user, err := h.userStore.GetUserByUserID(ctx, req.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewError("Failed to check user existence: "+err.Error()))
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, response.NewError("User not found"))
	}

	// Validation: Prevent modifying the admin user
	if user.RoleID == userrole.Admin {
		return c.JSON(http.StatusForbidden, response.NewError("Cannot modify an admin user"))
	}

	// Modify user's DisabledAt to a future date
	if err := h.userStore.ModifyUserValidity(ctx, req.UserID, time.Date(2100, 12, 31, 23, 59, 59, 0, time.UTC)); err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewError("Failed to activate user: "+err.Error()))
	}

	return c.JSON(http.StatusOK, response.NewSuccess("User activated successfully"))
}
