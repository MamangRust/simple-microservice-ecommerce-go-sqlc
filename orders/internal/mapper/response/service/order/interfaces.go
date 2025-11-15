package orderservicemapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
)

type OrderBaseResponseMapper interface {
	ToOrderResponse(Order *record.OrderRecord) *response.OrderResponse
}

type OrderQueryResponseMapper interface {
	OrderBaseResponseMapper

	ToOrdersResponse(Orders []*record.OrderRecord) []*response.OrderResponse

	ToOrdersResponseDeleteAt(Orders []*record.OrderRecord) []*response.OrderResponseDeleteAt
}

type OrderCommandResponseMapper interface {
	OrderBaseResponseMapper

	ToOrderResponseDeleteAt(Order *record.OrderRecord) *response.OrderResponseDeleteAt
}
