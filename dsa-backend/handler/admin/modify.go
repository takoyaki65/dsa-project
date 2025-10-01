package admin

import (
	"context"
	"dsa-backend/handler/response"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/takoyaki65/dsa-project/database/model/userrole"
	"golang.org/x/crypto/bcrypt"
)

type ModifyUserRequest struct {
	UserID   string  `param:"user_id" validate:"required"`
	UserName *string `json:"name" validate:"omitempty,min=1,max=30"`
	Password *string `json:"password" validate:"omitempty,min=8,max=100"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Role     *string `json:"role" validate:"omitempty,oneof=admin manager student"`
}

// ModifyUser modifies user details.
//
//	@Summary		Modify user details
//	@Description	Modify user details such as name, password, email, and role.
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			user_id	path		string				true	"User ID"
//	@Param			user	body		ModifyUserRequest	true	"User modification details"
//	@Success		200		{object}	response.Success	"User modified successfully"
//	@Failure		400		{object}	response.Error		"Invalid request body or validation failed"
//	@Failure		403		{object}	response.Error		"Cannot modify an admin user"
//	@Failure		404		{object}	response.Error		"User not found"
//	@Failure		500		{object}	response.Error		"Failed to modify user"
//	@Security		OAuth2Password[admin]
//	@Router			/admin/modify/{user_id} [patch]
func (h *Handler) ModifyUser(c echo.Context) error {
	ctx := context.Background()
	var req ModifyUserRequest
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

	var new_hashed_password *string = nil
	if req.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, response.NewError("Failed to hash password: "+err.Error()))
		}
		hashedPasswordStr := string(hashedPassword)
		new_hashed_password = &hashedPasswordStr
	}

	// convert role string to role id
	var new_role_id *userrole.Role = nil
	if req.Role != nil {
		roleID, ok := userrole.RoleNameToUserRole()[userrole.RoleName(*req.Role)]
		if !ok {
			return c.JSON(http.StatusBadRequest, response.NewError("Invalid role"))
		}
		new_role_id = &roleID
	}

	if err := h.userStore.ModifyUserDetails(ctx, req.UserID, req.UserName, new_hashed_password, req.Email, new_role_id); err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewError("Failed to modify user: "+err.Error()))
	}

	panic("not implemented")
}
