package rolehandler

import (
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/service"
	"github.com/MamangRust/simple_microservice_ecommerce/role/pkg/logger"
)

type RoleHandler interface {
	RoleQueryHandleGrpc
	RoleCommandHandleGrpc
}

type roleHandler struct {
	RoleQueryHandleGrpc
	RoleCommandHandleGrpc
}

func NewRoleHandler(service service.Service, logger logger.LoggerInterface) RoleHandler {
	return &roleHandler{
		RoleQueryHandleGrpc:   NewRoleQueryHandleGrpc(service, logger),
		RoleCommandHandleGrpc: NewRoleCommandHandleGrpc(service, logger),
	}
}
