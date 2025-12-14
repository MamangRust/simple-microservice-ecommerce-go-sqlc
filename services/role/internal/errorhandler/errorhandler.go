package errorhandler

import "github.com/MamangRust/simple_microservice_ecommerce/role/pkg/logger"

type ErrorHandler struct {
	RoleQueryError       RoleQueryErrorHandler
	RoleCommandError     RoleCommandErrorHandler
	UserRoleCommandError UserRoleErrorHandler
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		RoleQueryError:       NewRoleQueryError(logger),
		RoleCommandError:     NewRoleCommandError(logger),
		UserRoleCommandError: NewUserRoleError(logger),
	}
}
