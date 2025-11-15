package orderitemservicemapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
)

type orderItemResponseMapper struct{}

func NewOrderItemResponseMapper() OrderItemResponseMapper {
	return &orderItemResponseMapper{}
}

func (s *orderItemResponseMapper) ToOrderItemResponse(order *record.OrderItemRecord) *response.OrderItemResponse {
	return &response.OrderItemResponse{
		ID:        order.ID,
		OrderID:   order.OrderID,
		ProductID: order.ProductID,
		Quantity:  order.Quantity,
		Price:     order.Price,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}
}

func (s *orderItemResponseMapper) ToOrderItemsResponse(orders []*record.OrderItemRecord) []*response.OrderItemResponse {
	var responses []*response.OrderItemResponse

	for _, order := range orders {
		responses = append(responses, s.ToOrderItemResponse(order))
	}

	return responses
}

func (s *orderItemResponseMapper) ToOrderItemResponseDeleteAt(order *record.OrderItemRecord) *response.OrderItemResponseDeleteAt {
	return &response.OrderItemResponseDeleteAt{
		ID:        order.ID,
		OrderID:   order.OrderID,
		ProductID: order.ProductID,
		Quantity:  order.Quantity,
		Price:     order.Price,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
		DeletedAt: order.DeletedAt,
	}
}

func (s *orderItemResponseMapper) ToOrderItemsResponseDeleteAt(orders []*record.OrderItemRecord) []*response.OrderItemResponseDeleteAt {
	var responses []*response.OrderItemResponseDeleteAt

	for _, order := range orders {
		responses = append(responses, s.ToOrderItemResponseDeleteAt(order))
	}

	return responses
}
