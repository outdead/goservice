package system

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/outdead/goservice/internal/server/http/response"
)

// Handler is wrapper for HTTP API handle functions in health group.
type Handler struct{}

// NewHandler creates new Handler.
func NewHandler() *Handler {
	return &Handler{}
}

// Ping responses pong for `health/ping` HTTP request.
func (h *Handler) Ping() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, response.Response{
			Message: "pong",
		})
	}
}
