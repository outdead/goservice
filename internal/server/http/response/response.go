package response

import (
	"net/http"

	"github.com/labstack/echo"
)

// Response contains fields for JSON response from HTTP server.
type Response struct {
	Code    int         `json:"errorCode"`
	Message string      `json:"errorMessage"`
	Result  interface{} `json:"result"`
	Count   int         `json:"count,omitempty"`
}

// Fields allows you to enumerate the fields included in the Result block of
// the response structure.
type Fields map[string]interface{}

// ServeResult sends a JSON response with the result data and an additional
// count, which is responsible for the amount of data suitable for returning
// in pagination mode.
func ServeResult(c echo.Context, result interface{}, count ...int) error {
	response := Response{
		Code:    0,
		Message: "",
		Result:  result,
	}

	if len(count) != 0 {
		response.Count = count[0]
	}

	return c.JSON(http.StatusOK, response)
}

// ServeValidateError sends a JSON response with 422 code and the passed error
// message. The 422 response code reports a validation error.
func ServeValidateError(c echo.Context, msg ...string) error {
	return Serve(c, http.StatusUnprocessableEntity, msg...)
}

// ServeNotFoundError sends a JSON response with a 404 code and the passed error
// message. The 404 response code notifies of the absence of the requested data.
func ServeNotFoundError(c echo.Context, msg ...string) error {
	return Serve(c, http.StatusNotFound, msg...)
}

// ServeInternalServerError sends a JSON response with a 500 code and the passed
// error message. The 500 response code notifies an internal server error.
// As a rule, the error that occurred directly is not displayed, but hidden and
// displayed only in the service logs.
func ServeInternalServerError(c echo.Context, msg ...string) error {
	return Serve(c, http.StatusInternalServerError, msg...)
}

// Serve sends a JSON response with the passed code and error message.
// It is possible not to pass an error message - in this case it will be taken
// based on the response code.
func Serve(c echo.Context, code int, msg ...string) error {
	if len(msg) != 0 {
		return serve(c, code, msg[0])
	}

	return serve(c, code, http.StatusText(code))
}

func serve(c echo.Context, code int, message string) error {
	return c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Result:  []int{},
	})
}
