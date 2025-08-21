package handler

import (
	"context"
	"crypto/rand"
	"dsa-backend/storage"
	"dsa-backend/utils"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

type Handler struct {
	userStore    storage.UserStore
	problemStore storage.ProblemStore
	jwtSecret    string
}

func NewHandler(db *bun.DB) *Handler {
	return &Handler{
		userStore:    *storage.NewUserStore(db),
		problemStore: *storage.NewProblemStore(db),
		jwtSecret:    generateSecretKey(),
	}
}

func (h *Handler) RegisterRoutes(r *echo.Group) {
	r.POST("/login", h.Login)

	userRouter := r.Group("/user", JWTMiddleware(h.jwtSecret), h.CheckValidityOfJWTMiddleware())
	userRouter.GET("/me", h.GetCurrentUser, RequiredScopesMiddleware("me"))

	problemRouter := r.Group("/problem", JWTMiddleware(h.jwtSecret), h.CheckValidityOfJWTMiddleware())
	// fetch problem info
	problemRouter.GET("/list", h.ListProblems)
	problemRouter.GET("/detail/:lectureid/:problemid", h.GetProblemInfo)
	// submit validation / judge request
	problemRouter.POST("/submit/validate/:lectureid/:problemid", h.ValidateSubmission)
	problemRouter.POST("/submit/judge/:lectureid/:problemid", h.JudgeSubmission, RequiredScopesMiddleware("grading"))

	problemRouter.PUT("/create", h.CreateLectureEntry, RequiredScopesMiddleware("grading"))
	problemRouter.PATCH("/update/:lectureid", h.UpdateLectureEntry, RequiredScopesMiddleware("grading"))
	problemRouter.DELETE("/delete/:lectureid", h.DeleteLectureEntry, RequiredScopesMiddleware("grading"))
	problemRouter.GET("/create/:lectureid", h.RegisterProblem, RequiredScopesMiddleware("grading"))
	problemRouter.DELETE("/delete/:lectureid/:problemid", h.DeleteProblem, RequiredScopesMiddleware("grading"))

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

func (h *Handler) CheckValidityOfJWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, ok := c.Get("user").(*jwt.Token)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			claims, ok := token.Claims.(*utils.JwtCustomClaims)

			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
			}

			// Check token expiration
			if claims.ExpiresAt.Time.Before(time.Now()) {
				return echo.NewHTTPError(http.StatusUnauthorized, "token expired")
			}

			ctx := context.Background()
			// Check login history existence
			loginHistory, err := h.userStore.GetLoginHistory(&ctx, claims.UserID, claims.IssuedAt.Time)

			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to get login history")
			}

			if loginHistory == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "login history not found")
			}

			return next(c)
		}
	}
}
