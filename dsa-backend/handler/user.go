package handler

import (
	"context"
	"dsa-backend/storage/model"
	"dsa-backend/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Login godoc
//
//	@Summary		User Login
//	@Description	User login with user ID and password. Returns a JWT token if successful.
//	@Tags			user
//	@Accept			x-www-form-urlencoded
//	@Product		json
//	@param			username	formData	string				true	"User ID"
//	@param			password	formData	string				true	"Password"
//	@Success		200			{object}	userLoginResponse	"Login successful. Returns a JWT token."
//	@Failure		400			{object}	ErrorResponse		"Bad request. This error occurs if the user ID or password is missing or incorrect."
//	@Failure		500			{string}	ErrorResponse		"Internal server error. This error occurs if there is an issue with the database or password hashing."
//	@Router			/login [post]
func (h *Handler) Login(c echo.Context) error {
	ctx := context.Background()

	var loginRequest userLoginRequest
	err := loginRequest.bind(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newErrorResponse("failed to bind request: "+err.Error()))
	}

	plain_password := loginRequest.Password

	userRecord, err := h.userStore.GetUserByUserID(&ctx, loginRequest.UserId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, newErrorResponse("failed to get user: "+err.Error()))
	}

	if userRecord == nil {
		return c.JSON(http.StatusBadRequest, newErrorResponse("user not found"))
	}

	hashed_password := userRecord.HashedPassword

	// Verify provided password against the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(hashed_password), []byte(plain_password))

	if err != nil {
		return c.JSON(http.StatusBadRequest, newErrorResponse("wrong userid or password"))
	}

	// get user role
	userRole := userRecord.UserRole.Name
	// get user scopes
	scopes, err := GetUserScopes(userRole)

	if err != nil {
		return c.String(http.StatusInternalServerError, "invalid user role: "+userRole)
	}

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(time.Hour * 12) // 12 hours expiration

	// register LoginHistory
	{
		err := h.userStore.RegisterLoginHistory(&ctx, &model.LoginHistory{
			UserID:  userRecord.UserID,
			LoginAt: issuedAt,
		})

		if err != nil {
			return c.JSON(http.StatusInternalServerError, newErrorResponse("failed to register login history: "+err.Error()))
		}
	}

	// create JWT token
	token, err := utils.IssueNewToken(userRecord.UserID, scopes, h.jwtSecret, issuedAt, expiredAt)

	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to issue token")
	}

	return c.JSON(http.StatusOK, userLoginResponse{
		Token:     token,
		TokenType: "bearer",
		ExpiredAt: expiredAt.Unix(),
		User: userResponse{
			ID:    userRecord.UserID,
			Name:  userRecord.Name,
			Email: userRecord.Email,
		},
	})
}

// GetCurrentUser godoc
//
//	@Summary		Get current user information
//	@Description	Get current user information from JWT token
//	@Tags			user
//	@Product		json
//	@Success		200	{object}	userResponse	"Current user information"
//	@Failure		401	{object}	ErrorResponse	"Unauthorized"
//	@Failure		500	{string}	ErrorResponse	"Internal server error"
//	@Security		OAuth2Password[me]
//	@Router			/user/me [get]
func (h *Handler) GetCurrentUser(c echo.Context) error {
	ctx := context.Background()
	// Get userID from jwt token
	claim, err := utils.GetJWTClaims(&c)
	if err != nil {
		return err
	}

	// Get User data from db
	userRecord, err := h.userStore.GetUserByUserID(&ctx, claim.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newErrorResponse("failed to get user: "+err.Error()))
	}

	if userRecord == nil {
		return c.JSON(http.StatusUnauthorized, newErrorResponse("user not found"))
	}

	return c.JSON(http.StatusOK, userResponse{
		ID:    userRecord.UserID,
		Name:  userRecord.Name,
		Email: userRecord.Email,
	})
}
