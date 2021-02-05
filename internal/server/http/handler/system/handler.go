package system

import (
	"github.com/labstack/echo/v4"
	"github.com/outdead/goservice/internal/server/http/middleware"
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

// ContextTest responses pong for context HTTP request. Adds custom text
// from context.
func (h *Handler) ContextTest() echo.HandlerFunc {
	return middleware.HandlerFunc(func(ctx *middleware.Context) error {
		return ctx.ServeResult(map[string]interface{}{
			"data":       "pong",
			"request_id": ctx.RequestID,
		})
	})
}
