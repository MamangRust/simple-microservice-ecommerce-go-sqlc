package errorhandler

import (
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/response"
	userroleserviceerrors "github.com/MamangRust/simple_microservice_ecommerce/role/pkg/errors/user_role_errors/service"
	"github.com/MamangRust/simple_microservice_ecommerce/role/pkg/logger"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type userRoleCommandError struct {
	logger logger.LoggerInterface
}

func NewUserRoleError(logger logger.LoggerInterface) *userRoleCommandError {
	return &userRoleCommandError{
		logger: logger,
	}
}

func (e *userRoleCommandError) HandleAssignRoleToUserError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.UserRoleResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserRoleResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		userroleserviceerrors.ErrFailedAssignRoleToUser,
		fields...,
	)
}

func (e *userRoleCommandError) HandleUpdateRoleToUserError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.UserRoleResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserRoleResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		userroleserviceerrors.ErrFailedUpdateRoleToUser,
		fields...,
	)
}
