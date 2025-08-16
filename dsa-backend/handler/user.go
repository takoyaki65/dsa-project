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

func (h *Handler) Login(c echo.Context) error {
	ctx := context.Background()
	var loginRequest userLoginRequest
	err := c.Bind(&loginRequest)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	plain_password := loginRequest.User.Password

	userRecord, err := h.userStore.GetUserByUserID(&ctx, loginRequest.User.UserId)

	if err != nil {
		return c.String(http.StatusBadRequest, "wrong userid or password")
	}

	hashed_password := userRecord.HashedPassword

	// Authenticate user
	err = bcrypt.CompareHashAndPassword([]byte(hashed_password), []byte(plain_password))

	if err != nil {
		return c.String(http.StatusBadRequest, "wrong userid or password")
	}

	// get user role
	userRole := userRecord.UserRole.Name
	// get user scopes
	scopes, err := GetUserScopes(userRole)

	if err != nil {
		return c.String(http.StatusInternalServerError, "invalid user role: "+userRole)
	}

	// create JWT token
	token, err := utils.IssueNewToken(userRecord.UserID, scopes, h.jwtSecret)

	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to issue token")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}

func (h *Handler) CreateAdminUser(c echo.Context) error {
	userid := c.FormValue("userid")
	username := c.FormValue("username")
	pasword := c.FormValue("password")
	email := c.FormValue("email")
	if userid == "" || username == "" || pasword == "" {
		return c.String(http.StatusBadRequest, "userid, username, and password are required")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pasword), bcrypt.DefaultCost)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to hash password")
	}

	adminID, err := GetRoleID("admin")
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to get admin role ID")
	}

	// Create the user
	ctx := context.Background()
	// TODO: Implement the logic to create an admin user in the database
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
		return c.String(http.StatusInternalServerError, "failed to create admin user: "+err.Error())
	}

	go func() {
		time.Sleep(2 * time.Second)
		close(h.shutdownChan)
	}()

	return c.String(http.StatusOK, "admin user created successfully. Server will shutdown for restart.")
}
