package errorhandler

import (
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/response"
	roleserviceerrors "github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/errors/role/service"
	userserviceerrors "github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/errors/user/service"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/logger"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type registerError struct {
	logger logger.LoggerInterface
}

func NewRegisterError(logger logger.LoggerInterface) *registerError {
	return &registerError{
		logger: logger,
	}
}

func (e *registerError) HandleAssignRoleError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		userserviceerrors.ErrUserNotFoundRes,
		fields...,
	)
}

func (e *registerError) HandleFindEmailError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		userserviceerrors.ErrUserNotFoundRes,
		fields...,
	)
}

func (e *registerError) HandleFindRoleError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		roleserviceerrors.ErrRoleNotFoundRes,
		fields...,
	)
}

func (e *registerError) HandleCreateUserError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		userserviceerrors.ErrUserNotFoundRes,
		fields...,
	)
}
