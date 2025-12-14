package middlewares

import (
	"net/http"
	"reflect"

	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/domain/response"
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/pkg/logger"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCErrorHandlingMiddleware interface {
	HandleGRPCError(c echo.Context, err error, serviceName string) error
	ValidateClient(client interface{}, serviceName string) error
	ServiceUnavailableMiddleware(serviceName string) echo.MiddlewareFunc
}

type gRPCErrorHandlingMiddleware struct {
	logger logger.LoggerInterface
}

func NewGRPCErrorHandlingMiddleware(logger logger.LoggerInterface) GRPCErrorHandlingMiddleware {
	return &gRPCErrorHandlingMiddleware{
		logger: logger,
	}
}

func (m *gRPCErrorHandlingMiddleware) HandleGRPCError(c echo.Context, err error, serviceName string) error {
	if err == nil {
		return nil
	}

	if grpcErr, ok := status.FromError(err); ok {

		m.logger.Warn("gRPC error",
			zap.String("service", serviceName),
			zap.String("grpc_code", grpcErr.Code().String()),
			zap.String("grpc_message", grpcErr.Message()))

		switch grpcErr.Code() {

		case codes.NotFound:
			return jsonError(c, http.StatusNotFound, "resource not found", 404)

		case codes.Unauthenticated:
			return jsonError(c, http.StatusUnauthorized, "unauthorized", 401)

		case codes.PermissionDenied:
			return jsonError(c, http.StatusForbidden, "forbidden", 403)

		case codes.InvalidArgument:
			return jsonError(c, http.StatusBadRequest, "invalid request", 400)

		case codes.DeadlineExceeded:
			return jsonError(c, http.StatusRequestTimeout, "request timeout", 408)

		case codes.Unavailable:
			return jsonError(c, http.StatusServiceUnavailable, "service temporarily unavailable", 503)

		default:
			return jsonError(c, http.StatusInternalServerError, "internal server error", 500)
		}
	}

	m.logger.Error("Unexpected error",
		zap.String("service", serviceName),
		zap.Error(err))

	return jsonError(c, http.StatusInternalServerError, "internal server error", 500)
}

func (m *gRPCErrorHandlingMiddleware) ValidateClient(client interface{}, serviceName string) error {
	if client == nil || (reflect.ValueOf(client).Kind() == reflect.Ptr && reflect.ValueOf(client).IsNil()) {
		m.logger.Error("gRPC client is nil",
			zap.String("service", serviceName))
		return echo.NewHTTPError(http.StatusServiceUnavailable, map[string]string{
			"error": serviceName + " service unavailable",
			"code":  "SERVICE_UNAVAILABLE",
		})
	}
	return nil
}

func (m *gRPCErrorHandlingMiddleware) ServiceUnavailableMiddleware(serviceName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}
}

func jsonError(c echo.Context, httpStatus int, message string, code int) error {
	return c.JSON(httpStatus, &response.ErrorResponse{
		Status:  "error",
		Message: message,
		Code:    code,
	})
}
