package errorhandler

import (
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/response"
	userserviceerrors "github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/errors/user/service"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/logger"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type marshalError struct {
	logger logger.LoggerInterface
}

func NewMarshalError(logger logger.LoggerInterface) *marshalError {
	return &marshalError{
		logger: logger,
	}
}

func (e *marshalError) HandleMarshalRegisterError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorJSONMarshal[*response.UserResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		userserviceerrors.ErrFailedSendEmail,
		fields...,
	)
}

func (e *marshalError) HandleMarsalForgotPassword(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorJSONMarshal[bool](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		userserviceerrors.ErrFailedSendEmail,
		fields...,
	)
}

func (e *marshalError) HandleMarshalVerifyCode(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorJSONMarshal[bool](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		userserviceerrors.ErrFailedSendEmail,
		fields...,
	)
}
