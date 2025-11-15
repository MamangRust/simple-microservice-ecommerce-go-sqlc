package handler

import (
	orderhandler "github.com/MamangRust/simple_microservice_ecommerce/order/internal/handler/order"
	orderitemhandler "github.com/MamangRust/simple_microservice_ecommerce/order/internal/handler/orderitem"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/service"
	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/logger"
)

type Handler interface {
	OrderHandler() orderhandler.OrderHandlerGrpc
	OrderItemHandler() orderitemhandler.OrderItemHandleGrpc
}

type handler struct {
	orderHandler     orderhandler.OrderHandlerGrpc
	orderItemHandler orderitemhandler.OrderItemHandleGrpc
}

func (h *handler) OrderHandler() orderhandler.OrderHandlerGrpc {
	return h.orderHandler
}
func (h *handler) OrderItemHandler() orderitemhandler.OrderItemHandleGrpc {
	return h.orderItemHandler
}

func NewHandler(service service.Service, logger logger.LoggerInterface) Handler {
	return &handler{
		orderHandler:     orderhandler.NewOrderHandleGrpc(service.OrderService(), logger),
		orderItemHandler: orderitemhandler.NewOrderItemHandleGrpc(service.OrderItemService(), logger),
	}
}
