package response

import (
	"net/http"

	"github.com/labstack/echo"
)

// Response содержит поля для JSON ответа от HTTP сервера.
type Response struct {
	Code    int         `json:"errorCode"`
	Message string      `json:"errorMessage"`
	Result  interface{} `json:"result"`
	Count   int         `json:"count,omitempty"`
}

// Fields позволяет перечислять поля, входящие в блок Result структуры ответа.
type Fields map[string]interface{}

// ServeResult отправляет JSON ответ с данными result и дополнительным послем
// count, отвечающим за количество данных, подходящих для возврата в режиме
// пагинации.
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

// ServeValidateError отправляет JSON ответ с кодом 422 и переданным сообщением
// об ошибке. Код ответа 422 уведомляет об ошибке валидации.
func ServeValidateError(c echo.Context, msg ...string) error {
	return Serve(c, http.StatusUnprocessableEntity, msg...)
}

// ServeNotFoundError отправляет JSON ответ с кодом 404 и переданным сообщением
// об ошибке. Код ответа 404 уведомляет об отсутсвии зпрашиваемых данных.
func ServeNotFoundError(c echo.Context, msg ...string) error {
	return Serve(c, http.StatusNotFound, msg...)
}

// ServeInternalServerError отправляет JSON ответ с кодом 500 и переданным
// сообщением об ошибке. Код ответа 500 уведомляет о внутренней ошибке сервера.
// Как правило непосредственно произошедшая ошибка не выводится, а скривается и
// отображается только в логах сервиса.
func ServeInternalServerError(c echo.Context, msg ...string) error {
	return Serve(c, http.StatusInternalServerError, msg...)
}

// Serve отправляет JSON ответ с переданным кодом и сообщением об ошибке.
// Можно не передавать сообщение об ошибке, в этом случе оно будет взято
// на основании кода ответа.
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
