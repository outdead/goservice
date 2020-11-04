package health

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/outdead/echo-skeleton/internal/api/httpserver/response"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Ping() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, response.Response{
			Message: "pong",
		})
	}
}
