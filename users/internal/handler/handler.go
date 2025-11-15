package handler

import (
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/service"
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/logger"
)

type Deps struct {
	Service service.Service
	Logger  logger.LoggerInterface
}

type Handler interface {
	UserQueryHandleGrpc
	UserCommandHandleGrpc
}

type handler struct {
	UserQueryHandleGrpc
	UserCommandHandleGrpc
}

func NewHandler(deps *Deps) Handler {
	return &handler{
		UserQueryHandleGrpc:   NewUserQueryHandleGrpc(deps.Service, deps.Logger),
		UserCommandHandleGrpc: NewUserCommandHandleGrpc(deps.Service, deps.Logger),
	}
}
