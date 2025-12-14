package errorhandler

import (
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/response"
	userserviceerrors "github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/errors/user/service"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/logger"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type loginError struct {
	logger logger.LoggerInterface
}

func NewLoginError(logger logger.LoggerInterface) *loginError {
	return &loginError{
		logger: logger,
	}
}

func (e *loginError) HandleFindEmailError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TokenResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TokenResponse](
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
