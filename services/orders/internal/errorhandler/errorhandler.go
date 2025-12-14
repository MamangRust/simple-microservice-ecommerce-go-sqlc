package errorhandler

import "github.com/MamangRust/simple_microservice_ecommerce/order/pkg/logger"

type ErrorHandler struct {
	OrderCommandError   OrderCommandError
	OrderQueryError     OrderQueryError
	OrderItemQueryError OrderItemQueryError
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		OrderCommandError:   NewOrderCommandError(logger),
		OrderQueryError:     NewOrderQueryError(logger),
		OrderItemQueryError: NewOrderItemQueryError(logger),
	}
}
