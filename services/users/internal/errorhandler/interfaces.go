package errorhandler

import (
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type UserCommandErrorHandler interface {
	HandleCreateUserError(
		err error, method, tracePrefix string, span trace.Span, status *string,
		fields ...zap.Field,
	) (*response.UserResponse, *response.ErrorResponse)

	HandleUpdateUserError(
		err error, method, tracePrefix string, span trace.Span, status *string,
		fields ...zap.Field,
	) (*response.UserResponse, *response.ErrorResponse)

	HandleTrashedUserError(
		err error, method, tracePrefix string, span trace.Span, status *string,
		fields ...zap.Field,
	) (*response.UserResponseDeleteAt, *response.ErrorResponse)

	HandleRestoreUserError(
		err error, method, tracePrefix string, span trace.Span, status *string,
		fields ...zap.Field,
	) (*response.UserResponseDeleteAt, *response.ErrorResponse)

	HandleDeleteUserPermanentError(
		err error, method, tracePrefix string, span trace.Span, status *string,
		fields ...zap.Field,
	) (bool, *response.ErrorResponse)

	HandleDeleteAllUserPermanentError(
		err error, method, tracePrefix string, span trace.Span, status *string,
		fields ...zap.Field,
	) (bool, *response.ErrorResponse)

	HandleRestoreAllUserError(
		err error, method, tracePrefix string, span trace.Span, status *string,
		fields ...zap.Field,
	) (bool, *response.ErrorResponse)
}

type UserQueryErrorHandler interface {
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.UserResponse, *int, *response.ErrorResponse)

	HandleRepositoryPaginationDeletedError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse)

	HandleRepositoryListError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.UserResponse, *response.ErrorResponse)

	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		defaultErr *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.UserResponse, *response.ErrorResponse)

	HandleRepositorySingleWithPasswordError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		defaultErr *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.UserResponseWithPassword, *response.ErrorResponse)
}
