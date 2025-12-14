package orderservicemapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
)

type orderCommandResponseMapper struct{}

func NewOrderCommandResponseMapper() OrderCommandResponseMapper {
	return &orderCommandResponseMapper{}
}

func (s *orderCommandResponseMapper) ToOrderResponse(order *record.OrderRecord) *response.OrderResponse {
	return &response.OrderResponse{
		ID:         order.ID,
		UserID:     order.UserID,
		TotalPrice: order.TotalPrice,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}
}

func (s *orderCommandResponseMapper) ToOrderResponseDeleteAt(order *record.OrderRecord) *response.OrderResponseDeleteAt {
	return &response.OrderResponseDeleteAt{
		ID:         order.ID,
		UserID:     order.UserID,
		TotalPrice: order.TotalPrice,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
		DeletedAt:  order.DeletedAt,
	}
}
