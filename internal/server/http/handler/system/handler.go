package system

import (
	"github.com/labstack/echo/v4"
	"github.com/outdead/goservice/internal/server/http/response"
)

// Handler is wrapper for HTTP API handle functions in health group.
type Handler struct{}

// NewHandler creates new Handler.
func NewHandler() *Handler {
	return &Handler{}
}

// Ping godoc
// @Summary System info
// @Description Get system info
// @Tags system
// @Accept  json
// @Produce  json
// @Success 200 {object} response.Response
// @Failure 200 {object} response.Response
// @Router /system/ping [get]
//
// Ping responses pong for ping HTTP request.
func (h *Handler) Ping(c echo.Context) error {
	return response.ServeResult(c, "pong")
}
