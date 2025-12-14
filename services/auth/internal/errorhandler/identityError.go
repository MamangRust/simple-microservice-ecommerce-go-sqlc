package errorhandler

import (
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/response"
	refreshtokenserviceerrors "github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/errors/refresh_token/service"
	userserviceerrors "github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/errors/user/service"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/logger"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type identityError struct {
	logger logger.LoggerInterface
}

func NewIdentityError(logger logger.LoggerInterface) *identityError {
	return &identityError{
		logger: logger,
	}
}

func (e *identityError) HandleInvalidTokenError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TokenResponse, *response.ErrorResponse) {
	return handleErrorTokenTemplate[*response.TokenResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		refreshtokenserviceerrors.ErrFailedInvalidToken,
		fields...,
	)
}

func (e *identityError) HandleExpiredRefreshTokenError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TokenResponse, *response.ErrorResponse) {
	return handleErrorTokenTemplate[*response.TokenResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		refreshtokenserviceerrors.ErrFailedParseExpirationDate,
		fields...,
	)
}

func HandleInvalidFormatUserIDError[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	return handleErrorInvalidID[T](
		logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		userserviceerrors.ErrUserNotFoundRes,
		fields...,
	)
}

func (e *identityError) HandleDeleteRefreshTokenError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TokenResponse, *response.ErrorResponse) {
	return handleErrorTokenTemplate[*response.TokenResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		refreshtokenserviceerrors.ErrFailedDeleteRefreshToken,
		fields...,
	)
}

func (e *identityError) HandleUpdateRefreshTokenError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TokenResponse, *response.ErrorResponse) {
	return handleErrorTokenTemplate[*response.TokenResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		refreshtokenserviceerrors.ErrFailedUpdateRefreshToken,
		fields...,
	)
}

func (e *identityError) HandleValidateTokenError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorTokenTemplate[*response.UserResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		refreshtokenserviceerrors.ErrFailedInvalidToken,
		fields...,
	)
}

func (e *identityError) HandleFindByIdError(
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

func (e *identityError) HandleGetMeError(
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
