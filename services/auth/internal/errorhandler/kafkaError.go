package errorhandler

import (
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/response"
	userserviceerrors "github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/errors/user/service"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/logger"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type kafkaError struct {
	logger logger.LoggerInterface
}

func NewKafkaError(logger logger.LoggerInterface) *kafkaError {
	return &kafkaError{
		logger: logger,
	}
}

func (e *kafkaError) HandleSendEmailForgotPassword(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorKafkaSend[bool](
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

func (e *kafkaError) HandleSendEmailRegister(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorKafkaSend[*response.UserResponse](
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

func (e *kafkaError) HandleSendEmailVerifyCode(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorKafkaSend[bool](
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
