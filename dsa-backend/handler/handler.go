package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

type Handler struct {
	db *bun.DB
}

func NewHandler(db *bun.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) RegisterRoutes(r *echo.Echo) {
	// not impelemnted yet
	panic("RegisterRoutes not implemented")
}
