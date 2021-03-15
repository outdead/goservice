package request

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
)

// ErrInvalidPathParam returned on parameter parsing errors.
var ErrInvalidPathParam = errors.New("invalid path param")

// PathParamInt gets a variable from part of the request path by name and
// converts it to int. If the variable was not passed, then an error is returned
// (formally this situation is impossible).
func PathParamInt(c echo.Context, name string) (int, error) {
	param := c.Param(name)
	if param == "" {
		return 0, fmt.Errorf("%s: %w int", name, ErrInvalidPathParam)
	}

	val, err := strconv.Atoi(param)
	if err != nil {
		return val, fmt.Errorf("%s: %w int", name, ErrInvalidPathParam)
	}

	return val, nil
}
