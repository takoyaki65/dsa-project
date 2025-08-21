package handler

import (
	"time"

	"github.com/labstack/echo/v4"
)

type userLoginRequest struct {
	UserId   string `form:"username" validate:"required"`
	Password string `form:"password" validate:"required"`
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

type userRegisterRequest struct {
	UserId   string `json:"userid" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"omitempty,email" default:"hoge@gmail.com"`
}

func (r *userRegisterRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	return nil
}

type LectureEntryRequest struct {
	ID        int64     `json:"id" validate:"required" default:"0"`
	Title     string    `json:"title" validate:"required"`
	StartDate time.Time `json:"start_date" validate:"required" default:"2025-10-01T10:00:00+09:00"`
	Deadline  time.Time `json:"deadline" validate:"required" default:"2025-12-01T10:00:00+09:00"`
}

func (le *LectureEntryRequest) bind(c echo.Context) error {
	if err := c.Bind(le); err != nil {
		return err
	}
	if err := c.Validate(le); err != nil {
		return err
	}
	return nil
}
