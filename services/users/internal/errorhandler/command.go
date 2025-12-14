package errorhandler

import (
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/response"
	userserviceerrors "github.com/MamangRust/simple_microservice_ecommerce/user/pkg/errors/user_errors/service"
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/logger"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type userCommandError struct {
	logger logger.LoggerInterface
}

func NewUserCommandError(logger logger.LoggerInterface) UserCommandErrorHandler {
	return &userCommandError{logger: logger}
}

func (e *userCommandError) HandleCreateUserError(
	err error, method, tracePrefix string,
	span trace.Span, status *string,
	fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserResponse](
		e.logger, err, method, tracePrefix, span, status,
		userserviceerrors.ErrFailedCreateUser,
		fields...,
	)
}

func (e *userCommandError) HandleUpdateUserError(
	err error, method, tracePrefix string,
	span trace.Span, status *string,
	fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserResponse](
		e.logger, err, method, tracePrefix, span, status,
		userserviceerrors.ErrFailedUpdateUser,
		fields...,
	)
}

func (e *userCommandError) HandleTrashedUserError(
	err error, method, tracePrefix string,
	span trace.Span, status *string,
	fields ...zap.Field,
) (*response.UserResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserResponseDeleteAt](
		e.logger, err, method, tracePrefix, span, status,
		userserviceerrors.ErrFailedTrashedUser,
		fields...,
	)
}

func (e *userCommandError) HandleRestoreUserError(
	err error, method, tracePrefix string,
	span trace.Span, status *string,
	fields ...zap.Field,
) (*response.UserResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserResponseDeleteAt](
		e.logger, err, method, tracePrefix, span, status,
		userserviceerrors.ErrFailedRestoreUser,
		fields...,
	)
}

func (e *userCommandError) HandleDeleteUserPermanentError(
	err error, method, tracePrefix string,
	span trace.Span, status *string,
	fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorBool(
		e.logger, err, method, tracePrefix, span, status,
		userserviceerrors.ErrFailedDeletePermanent,
		fields...,
	)
}

func (e *userCommandError) HandleDeleteAllUserPermanentError(
	err error, method, tracePrefix string,
	span trace.Span, status *string,
	fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorBool(
		e.logger, err, method, tracePrefix, span, status,
		userserviceerrors.ErrFailedDeleteAll,
		fields...,
	)
}

func (e *userCommandError) HandleRestoreAllUserError(
	err error, method, tracePrefix string,
	span trace.Span, status *string,
	fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorBool(
		e.logger, err, method, tracePrefix, span, status,
		userserviceerrors.ErrFailedRestoreAll,
		fields...,
	)
}
