package middleware

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Recover returns a middleware which recovers from panics anywhere in the chain
// and handles the control to the centralized HTTPErrorHandler.
func Recover() echo.MiddlewareFunc {
	return middleware.Recover()
}
