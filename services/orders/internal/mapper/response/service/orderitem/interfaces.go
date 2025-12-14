package orderitemservicemapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
)

type OrderItemResponseMapper interface {
	ToOrderItemResponse(order *record.OrderItemRecord) *response.OrderItemResponse

	ToOrderItemsResponse(orders []*record.OrderItemRecord) []*response.OrderItemResponse
	ToOrderItemResponseDeleteAt(order *record.OrderItemRecord) *response.OrderItemResponseDeleteAt

	ToOrderItemsResponseDeleteAt(orders []*record.OrderItemRecord) []*response.OrderItemResponseDeleteAt
}
