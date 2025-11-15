package orderhandler

import (
	orderservice "github.com/MamangRust/simple_microservice_ecommerce/order/internal/service/order"
	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/logger"
)

type OrderHandlerGrpc interface {
	OrderQueryHandleGrpc
	OrderCommandHandleGrpc
}

type orderHandler struct {
	OrderQueryHandleGrpc
	OrderCommandHandleGrpc
}

func NewOrderHandleGrpc(service orderservice.OrderService, logger logger.LoggerInterface) OrderHandlerGrpc {
	return &orderHandler{
		OrderQueryHandleGrpc:   NewOrderQueryHandleGrpc(service, logger),
		OrderCommandHandleGrpc: NewOrderCommandHandleGrpc(service, logger),
	}
}
