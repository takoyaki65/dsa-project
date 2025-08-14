package handler

import "github.com/labstack/echo/v4"

type userLoginRequest struct {
	User struct {
		UserId   string `json:"userid" validate:"required"`
		Password string `json:"password" validate:"required"`
	} `json:"user"`
}

func (r *userLoginRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	return nil
}
