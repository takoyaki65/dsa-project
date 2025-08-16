package handler

import (
	"context"
	"dsa-backend/utils"
	"net/http"

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
