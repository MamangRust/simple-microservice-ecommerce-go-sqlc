package response

import (
	"github.com/labstack/echo/v4"
)

func NewErrorResponse(message string, code int) *ErrorResponse {
	return &ErrorResponse{
		Status:  "error",
		Message: message,
		Code:    code,
	}
}

func NewApiErrorResponse(c echo.Context, statusText string, message string, code int) error {
	return c.JSON(code, ErrorResponse{
		Status:  statusText,
		Message: message,
		Code:    code,
	})
}
