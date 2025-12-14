package handler

import (
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/service"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/logger"
)

type Deps struct {
	Service service.Service
	Logger  logger.LoggerInterface
}

type Handler interface {
	AuthHandleGrpc
}

type handler struct {
	AuthHandleGrpc
}

func NewHandler(deps *Deps) Handler {
	return &handler{
		AuthHandleGrpc: NewAuthHandleGrpc(
			deps.Service,
			deps.Logger,
		),
	}
}
