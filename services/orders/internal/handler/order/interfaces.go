package orderhandler

import pborder "github.com/MamangRust/simple_microservice_ecommerce_pb/order"

type OrderQueryHandleGrpc interface {
	pborder.OrderQueryServiceServer
}

type OrderCommandHandleGrpc interface {
	pborder.OrderCommandServiceServer
}
