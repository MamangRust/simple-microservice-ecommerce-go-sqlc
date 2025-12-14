package errorhandler

import (
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/response"
	userserviceerrors "github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/errors/user/service"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/logger"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type passwordResetError struct {
	logger logger.LoggerInterface
}

func NewPasswordResetError(logger logger.LoggerInterface) *passwordResetError {
	return &passwordResetError{
		logger: logger,
	}
}

func (e *passwordResetError) HandleFindEmailError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
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

func (e *passwordResetError) HandleCreateResetTokenError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		userserviceerrors.ErrInternalServerError,
		fields...,
	)
}

func (e *passwordResetError) HandleFindTokenError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
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

func (e *passwordResetError) HandleUpdatePasswordError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
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

func (e *passwordResetError) HandleDeleteTokenError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
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

func (e *passwordResetError) HandleUpdateVerifiedError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
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

func (e *passwordResetError) HandleVerifyCodeError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
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
