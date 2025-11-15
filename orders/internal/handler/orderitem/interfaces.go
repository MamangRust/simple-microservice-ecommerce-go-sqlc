package orderitemhandler

import pborderitem "github.com/MamangRust/simple_microservice_ecommerce_pb/order_item"

type OrderItemHandleGrpc interface {
	pborderitem.OrderItemServiceServer
}
