package handler

import (
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/service"
	"github.com/MamangRust/simple_microservice_ecommerce/product/pkg/logger"
)

type Handler interface {
	ProductQueryHandleGrpc
	ProductCommandHandleGrpc
}

type handler struct {
	ProductQueryHandleGrpc
	ProductCommandHandleGrpc
}

func NewHandler(Service service.Service, logger logger.LoggerInterface) Handler {
	return &handler{
		ProductQueryHandleGrpc:   NewProductQueryHandleGrpc(Service, logger),
		ProductCommandHandleGrpc: NewProductCommandHandleGrpc(Service, logger),
	}
}
