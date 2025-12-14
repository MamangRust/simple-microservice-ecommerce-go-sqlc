package errorhandler

import (
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/response"
	refreshtokenserviceerrors "github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/errors/refresh_token/service"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/logger"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type tokenError struct {
	logger logger.LoggerInterface
}

func NewTokenError(logger logger.LoggerInterface) *tokenError {
	return &tokenError{
		logger: logger,
	}
}

func (e *tokenError) HandleCreateAccessTokenError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TokenResponse, *response.ErrorResponse) {
	return handleErrorTokenTemplate[*response.TokenResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		refreshtokenserviceerrors.ErrFailedCreateAccess,
		fields...,
	)
}

func (e *tokenError) HandleCreateRefreshTokenError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TokenResponse, *response.ErrorResponse) {
	return handleErrorTokenTemplate[*response.TokenResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		refreshtokenserviceerrors.ErrFailedCreateRefresh,
		fields...,
	)
}
