package router

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func New() *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS setting
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // Allow all origins (TODO: we need to restrict this in production)
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.POST, echo.DELETE, echo.PATCH, echo.OPTIONS},
	}))

	// Set custom validator for URL path parameter, query parameter, header and request body.
	e.Validator = NewValidator()

	// JWT auth settings
	e.Use(echojwt.WithConfig(echojwt.Config{
		// We create a new JWT signing key everytime we start the server.
		// So previous tokens will be invalidated
		SigningKey: []byte(generateSecretKey()),
		// User information extracted from JWT token can be referenced using "token" key.
		ContextKey: "token",
	}))

	return e
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
