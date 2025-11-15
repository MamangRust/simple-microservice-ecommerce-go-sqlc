package middlewares

import (
	"net/http"

	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/domain/response"
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/pkg/logger"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ErrorHandlerMiddleware struct {
	logger logger.LoggerInterface
}

func NewErrorHandlerMiddleware(logger logger.LoggerInterface) *ErrorHandlerMiddleware {
	return &ErrorHandlerMiddleware{
		logger: logger,
	}
}

func (m *ErrorHandlerMiddleware) ErrorHandler(err error, c echo.Context) {
	if err == nil {
		return
	}

	m.logger.Error("Request error occurred",
		zap.String("method", c.Request().Method),
		zap.String("path", c.Request().URL.Path),
		zap.String("remote_ip", c.RealIP()),
		zap.Error(err),
	)

	if httpErr, ok := err.(*echo.HTTPError); ok {
		m.handleHTTPError(c, httpErr)
		return
	}

	c.JSON(http.StatusInternalServerError, map[string]string{
		"error": "internal server error",
		"code":  "INTERNAL_ERROR",
	})
}

func (m *ErrorHandlerMiddleware) handleHTTPError(c echo.Context, err *echo.HTTPError) {
	code := err.Code
	message := "internal server error"
	errorCode := "INTERNAL_ERROR"

	switch code {
	case http.StatusBadRequest:
		message = "bad request"
		errorCode = "BAD_REQUEST"
	case http.StatusUnauthorized:
		message = "unauthorized"
		errorCode = "UNAUTHORIZED"
	case http.StatusForbidden:
		message = "forbidden"
		errorCode = "FORBIDDEN"
	case http.StatusNotFound:
		message = "not found"
		errorCode = "NOT_FOUND"
	case http.StatusMethodNotAllowed:
		message = "method not allowed"
		errorCode = "METHOD_NOT_ALLOWED"
	case http.StatusRequestTimeout:
		message = "request timeout"
		errorCode = "TIMEOUT"
	case http.StatusTooManyRequests:
		message = "too many requests"
		errorCode = "TOO_MANY_REQUESTS"
	case http.StatusInternalServerError:
		message = "internal server error"
		errorCode = "INTERNAL_ERROR"
	case http.StatusServiceUnavailable:
		message = "service unavailable"
		errorCode = "SERVICE_UNAVAILABLE"
	}

	response := &response.ErrorResponse{
		Status:  errorCode,
		Message: message,
		Code:    err.Code,
	}

	c.JSON(code, response)
}
