package handler

import (
	rolehandler "github.com/MamangRust/simple_microservice_ecommerce/role/internal/handler/role"
	userrolehandler "github.com/MamangRust/simple_microservice_ecommerce/role/internal/handler/user_role"
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/service"
	"github.com/MamangRust/simple_microservice_ecommerce/role/pkg/logger"
)

type Handler interface {
	RoleHandler() rolehandler.RoleHandler
	UserRoleHandler() userrolehandler.UserRoleHandleGrpc
}

type handler struct {
	roleHandler     rolehandler.RoleHandler
	userRoleHandler userrolehandler.UserRoleHandleGrpc
}

func (h *handler) RoleHandler() rolehandler.RoleHandler {
	return h.roleHandler
}

func (h *handler) UserRoleHandler() userrolehandler.UserRoleHandleGrpc {
	return h.userRoleHandler
}

func NewHandler(service service.Service, logger logger.LoggerInterface) Handler {
	return &handler{
		roleHandler:     rolehandler.NewRoleHandler(service, logger),
		userRoleHandler: userrolehandler.NewUserRoleHandleGrpc(service, logger),
	}
}
