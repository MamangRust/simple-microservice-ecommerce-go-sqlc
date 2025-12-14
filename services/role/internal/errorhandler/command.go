package errorhandler

import (
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/response"
	roleserviceerrors "github.com/MamangRust/simple_microservice_ecommerce/role/pkg/errors/role_errors/service"
	"github.com/MamangRust/simple_microservice_ecommerce/role/pkg/logger"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type roleCommandError struct {
	logger logger.LoggerInterface
}

func NewRoleCommandError(logger logger.LoggerInterface) *roleCommandError {
	return &roleCommandError{
		logger: logger,
	}
}

func (e *roleCommandError) HandleCreateRoleError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.RoleResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.RoleResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		roleserviceerrors.ErrFailedCreateRole,
		fields...,
	)
}

func (e *roleCommandError) HandleUpdateRoleError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.RoleResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.RoleResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		roleserviceerrors.ErrFailedUpdateRole,
		fields...,
	)
}

func (e *roleCommandError) HandleTrashedRoleError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.RoleResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.RoleResponseDeleteAt](
		e.logger,
		err, method, tracePrefix, span, status,
		roleserviceerrors.ErrFailedTrashedRole,
		fields...,
	)
}

func (e *roleCommandError) HandleRestoreRoleError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.RoleResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.RoleResponseDeleteAt](
		e.logger,
		err, method, tracePrefix, span, status,
		roleserviceerrors.ErrFailedRestoreRole,
		fields...,
	)
}

func (e *roleCommandError) HandleDeleteRolePermanentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		roleserviceerrors.ErrFailedDeletePermanent,
		fields...,
	)
}

func (e *roleCommandError) HandleDeleteAllRolePermanentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		roleserviceerrors.ErrFailedDeleteAll,
		fields...,
	)
}

func (e *roleCommandError) HandleRestoreAllRoleError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		roleserviceerrors.ErrFailedRestoreAll,
		fields...,
	)
}
