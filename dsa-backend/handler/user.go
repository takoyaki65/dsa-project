package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Login(c echo.Context) error {
	var loginRequest userLoginRequest
	err := c.Bind(&loginRequest)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	panic("not implemented")
}
