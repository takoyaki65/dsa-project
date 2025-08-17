package handler

import (
	"context"
	"dsa-backend/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Login godoc
// @Summary User Login
// @Description User login with user ID and password. Returns a JWT token if successful.
// @Tags user
// @Accept json
// @Product json
// @Param user body userLoginRequest true "User login info"
// @Success 200 {object} userLoginResponse "Login successful. Returns a JWT token."
// @Failure 400 {object} utils.Error "Bad request. This error occurs if the user ID or password is missing or incorrect."
// @Failure 500 {string} utils.Error "Internal server error. This error occurs if there is an issue with the database or password hashing."
// @Router /login [post]
func (h *Handler) Login(c echo.Context) error {
	ctx := context.Background()
	var loginRequest userLoginRequest
	err := loginRequest.bind(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorWithMessage("failed to bind request: "+err.Error()))
	}

	plain_password := loginRequest.Password

	userRecord, err := h.userStore.GetUserByUserID(&ctx, loginRequest.UserId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorWithMessage("failed to get user: "+err.Error()))
	}

	if userRecord == nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorWithMessage("user not found"))
	}

	hashed_password := userRecord.HashedPassword

	// Verify provided password against the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(hashed_password), []byte(plain_password))

	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorWithMessage("wrong userid or password"))
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

	return c.JSON(http.StatusOK, userLoginResponse{
		Token: token,
	})
}
