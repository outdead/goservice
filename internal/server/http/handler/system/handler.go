package system

import (
	"github.com/labstack/echo"
	"github.com/outdead/goservice/internal/server/http/response"
)

// Handler is wrapper for HTTP API handle functions in health group.
type Handler struct{}

// NewHandler creates new Handler.
func NewHandler() *Handler {
	return &Handler{}
}

// Ping responses pong for ping HTTP request.
func (h *Handler) Ping(c echo.Context) error {
	return response.ServeResult(c, "pong")
}
