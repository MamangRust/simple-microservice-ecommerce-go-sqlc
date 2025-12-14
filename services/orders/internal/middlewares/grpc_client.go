package middlewares

import (
	"reflect"

	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCErrorHandling interface {
	HandleGRPCError(err error, serviceName string) *response.ErrorResponse
	ValidateClient(client interface{}, serviceName string) *response.ErrorResponse
}

type gRPCErrorHandling struct {
	logger logger.LoggerInterface
}

func NewGRPCErrorHandling(logger logger.LoggerInterface) GRPCErrorHandling {
	return &gRPCErrorHandling{
		logger: logger,
	}
}

func (m *gRPCErrorHandling) HandleGRPCError(err error, serviceName string) *response.ErrorResponse {
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
			return response.NewApiErrorResponse("error", "resource not found", 404)

		case codes.Unauthenticated:
			return response.NewApiErrorResponse("error", "unauthorized", 401)

		case codes.PermissionDenied:
			return response.NewApiErrorResponse("error", "forbidden", 403)

		case codes.InvalidArgument:
			return response.NewApiErrorResponse("error", "invalid request", 400)

		case codes.DeadlineExceeded:
			return response.NewApiErrorResponse("error", "request timeout", 408)

		case codes.Unavailable:
			return response.NewApiErrorResponse("error", "service temporarily unavailable", 503)

		default:
			return response.NewApiErrorResponse("error", "internal server error", 500)
		}
	}

	m.logger.Error("Unexpected non-gRPC error",
		zap.String("service", serviceName),
		zap.Error(err))

	return response.NewApiErrorResponse("error", "internal server error", 500)
}

func (m *gRPCErrorHandling) ValidateClient(client interface{}, serviceName string) *response.ErrorResponse {
	if client == nil || (reflect.ValueOf(client).Kind() == reflect.Ptr && reflect.ValueOf(client).IsNil()) {

		m.logger.Error("gRPC client is nil",
			zap.String("service", serviceName))

		return response.NewApiErrorResponse("error", serviceName+" service unavailable", 503)
	}

	return nil
}
