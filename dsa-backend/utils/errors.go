package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Error struct {
	Errors map[string]interface{} `json:"errors"`
}

func NewErrorWithMessage(message string) Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["body"] = message
	return e
}

func NewError(err error) Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	switch v := err.(type) {
	case *echo.HTTPError:
		e.Errors["body"] = v.Message
	default:
		e.Errors["body"] = v.Error()
	}
	return e
}

func NewValidationError(errs validator.ValidationErrors) Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	for _, v := range errs {
		e.Errors[v.Field()] = fmt.Sprintf("%v", v.Tag())
	}
	return e
}

func AccessForbidden() Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["body"] = "access forbidden"
	return e
}

func NotFound() Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["body"] = "not found"
	return e
}
