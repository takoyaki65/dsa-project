package handler

import (
	"context"
	"dsa-backend/model"
	"dsa-backend/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// CreateAdmin godoc
// @Summary Create Admin User. This endpoint is exposed only if no admin user exists.
// @Description Create an admin user with the provided credentials. Note that this endpoint is only available if no admin user exists.
// @Tags initialization
// @Accept json
// @Product json
// @Param user body userRegisterRequest true "User info for registration"
// @Success 200 {object} createUserSuccess "Admin user created successfully. Server will shutdown for restart."
// @Failure 400 {object} utils.Error "Bad request. This error occurs if the admin user already exists or if the required fields are missing."
// @Failure 500 {string} utils.Error "Internal server error. This error occurs if there is an issue with the database or password hashing."
// @Router /admin/create [post]
func (h *Handler) CreateAdminUser(c echo.Context) error {
	// Check if the admin user already exists
	{
		ctx := context.Background()
		admins, err := h.userStore.GetUserListByUserRole(&ctx, "admin")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewErrorWithMessage("failed to check admin user: "+err.Error()))
		}
		if len(*admins) > 0 {
			// If an admin user already exists, return a 400 Bad Request
			return c.JSON(http.StatusBadRequest, utils.NewErrorWithMessage("admin user already exists"))
		}
	}

	// Bind and validate the form data
	var registerRequest userRegisterRequest
	if err := registerRequest.bind(c); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorWithMessage("failed to bind request: "+err.Error()))
	}

	userid := registerRequest.UserId
	username := registerRequest.Username
	pasword := registerRequest.Password
	email := registerRequest.Email
	if userid == "" || username == "" || pasword == "" {
		return c.JSON(http.StatusBadRequest, utils.NewErrorWithMessage("userid, username, and password are required"))
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pasword), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorWithMessage("failed to hash password: "+err.Error()))
	}

	adminID, err := GetRoleID("admin")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorWithMessage("failed to get admin role ID: "+err.Error()))
	}

	// Create the user
	ctx := context.Background()
	user := &model.UserList{
		UserID:         userid,
		Name:           username,
		HashedPassword: string(hashedPassword),
		RoleID:         int64(adminID),
		DisabledAt:     time.Now().Add(2 * 365 * 24 * time.Hour),
		Email:          &email,
	}
	err = h.userStore.CreateUser(&ctx, user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorWithMessage("failed to create admin user: "+err.Error()))
	}

	go func() {
		time.Sleep(2 * time.Second)
		close(h.shutdownChan)
	}()

	return c.JSON(http.StatusOK, newCreateUserSuccess("admin user created successfully. Server will shutdown for restart."))
}
