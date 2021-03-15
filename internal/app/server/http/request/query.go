package request

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// ErrInvalidQueryParam returned on parsing variables from GET parameters errors.
var ErrInvalidQueryParam = errors.New("invalid query param")

// QueryParamInt gets a variable from GET request parameters and converts it
// to int. If the variable was not passed then 0 is returned.
func QueryParamInt(c echo.Context, name string) (value int, err error) {
	if param := c.QueryParam(name); param != "" {
		value, err = strconv.Atoi(param)
		if err != nil {
			return value, fmt.Errorf("%s: %w int", name, ErrInvalidQueryParam)
		}
	}

	return value, nil
}

// QueryParamInt64 gets a variable from GET request parameters and converts it
// to int64. If the variable was not passed then 0 is returned.
func QueryParamInt64(c echo.Context, name string) (value int64, err error) {
	if param := c.QueryParam(name); param != "" {
		value, err = strconv.ParseInt(param, 10, 64)
		if err != nil {
			return value, fmt.Errorf("%s: %w int64", name, ErrInvalidQueryParam)
		}
	}

	return value, nil
}

// QueryParamBool gets a variable from GET request parameters and converts it
// to bool. If the variable was not passed then false is returned.
func QueryParamBool(c echo.Context, name string) (value bool, err error) {
	if param := c.QueryParam(name); param != "" {
		switch param {
		case "true":
			return true, nil
		case "false":
			return false, nil
		default:
			return false, fmt.Errorf("%s: %w bool", name, ErrInvalidQueryParam)
		}
	}

	return false, nil
}

// QueryParamSlice gets a variable from GET request parameters and converts it
// to []string. If the variable was not passed then empty slice is returned.
func QueryParamSlice(c echo.Context, name string) []string {
	var values []string

	if params := c.QueryParam(name); params != "" {
		for _, param := range strings.Split(params, ",") {
			if param != "" {
				values = append(values, param)
			}
		}
	}

	return values
}

// QueryParamIntSlice gets a variable from GET request parameters and converts it
// to []int. If the variable was not passed then empty slice is returned.
func QueryParamIntSlice(c echo.Context, name string) (values []int, err error) {
	if params := c.QueryParam(name); params != "" {
		for _, param := range strings.Split(params, ",") {
			if param != "" {
				value, err := strconv.Atoi(param)
				if err != nil {
					return values, fmt.Errorf("%s: %w []int", name, ErrInvalidQueryParam)
				}

				values = append(values, value)
			}
		}
	}

	return values, nil
}

// QueryParamInt64Slice gets a variable from GET request parameters and converts it
// to []int64. If the variable was not passed then empty slice is returned.
func QueryParamInt64Slice(c echo.Context, name string) (values []int64, err error) {
	if params := c.QueryParam(name); params != "" {
		for _, param := range strings.Split(params, ",") {
			if param != "" {
				value, err := strconv.ParseInt(param, 10, 64)
				if err != nil {
					return values, fmt.Errorf("%s: %w []int64", name, ErrInvalidQueryParam)
				}

				values = append(values, value)
			}
		}
	}

	return values, nil
}

// QueryParamTime gets a string variable from GET request parameters and converts
// it to time.Time string layout. If the variable was not passed then time.Time{}
// is returned.
func QueryParamTime(c echo.Context, name, layout string) (value time.Time, err error) {
	if param := c.QueryParam(name); param != "" {
		value, err = time.Parse(layout, param)
		if err != nil {
			return value, fmt.Errorf("%s: %w time in layout %q", name, ErrInvalidQueryParam, layout)
		}
	}

	return value, nil
}

// QueryParamTime gets a timestamp variable from GET request parameters and
// converts it to time.Time string layout. If the variable was not passed
// then time.Time{} is returned.
func QueryParamTimeUnix(c echo.Context, name string) (value time.Time, err error) {
	unix, err := QueryParamInt64(c, name)
	if err != nil {
		return value, fmt.Errorf("%s: %w timestamp", name, ErrInvalidQueryParam)
	}

	// Created time from unix = 0 is not IsZero.
	// Therefore create time only from a non-zero value.
	if unix != 0 {
		return time.Unix(unix, 0), nil
	}

	return time.Time{}, nil
}
