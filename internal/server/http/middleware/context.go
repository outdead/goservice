package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// WithContext adds custom context to echo.
func WithContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(&Context{Context: c, RequestID: "123"})
		}
	}
}

// HandlerFunc converts echo context to custom context in handler functions.
func HandlerFunc(h func(ctx *Context) error) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return h(ctx.(*Context))
	}
}

// Context extends echo context.
type Context struct {
	echo.Context
	RequestID string
}

// ServeResult sends a JSON response with the result data and an additional
// count, which is responsible for the amount of data suitable for returning
// in pagination mode.
func (ctx *Context) ServeResult(result interface{}, count ...int) error {
	// Response contains fields for JSON response from HTTP server.
	type Response struct {
		Code    int         `json:"errorCode"`
		Message string      `json:"errorMessage"`
		Result  interface{} `json:"result"`
		Count   int         `json:"count,omitempty"`
	}

	response := Response{
		Code:    0,
		Message: "",
		Result:  result,
	}

	if len(count) != 0 {
		response.Count = count[0]
	}

	return ctx.JSON(http.StatusOK, response)
}
