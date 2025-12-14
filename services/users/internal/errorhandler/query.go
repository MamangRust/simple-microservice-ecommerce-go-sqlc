package errorhandler

import (
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/response"
	userserviceerrors "github.com/MamangRust/simple_microservice_ecommerce/user/pkg/errors/user_errors/service"
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/logger"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type userQueryError struct {
	logger logger.LoggerInterface
}

func NewUserQueryError(logger logger.LoggerInterface) UserQueryErrorHandler {
	return &userQueryError{logger: logger}
}

func (e *userQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.UserResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.UserResponse](
		e.logger, err, method, tracePrefix, span, status, userserviceerrors.ErrFailedFindAll, fields...,
	)
}

func (e *userQueryError) HandleRepositoryPaginationDeletedError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.UserResponseDeleteAt](e.
		logger, err, method, tracePrefix, span, status, errResp, fields...,
	)
}

func (e *userQueryError) HandleRepositoryListError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.UserResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.UserResponse](e.logger, err, method, tracePrefix, span, status, userserviceerrors.ErrFailedFindAll, fields...)
}

func (e *userQueryError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserResponse](e.logger, err, method, tracePrefix, span, status, defaultErr, fields...)
}

func (e *userQueryError) HandleRepositorySingleWithPasswordError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (*response.UserResponseWithPassword, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserResponseWithPassword](e.logger, err, method, tracePrefix, span, status, defaultErr, fields...)
}
