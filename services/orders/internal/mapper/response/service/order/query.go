package orderservicemapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
)

type orderQueryResponseMapper struct {
}

func NewOrderQueryResponseMapper() OrderQueryResponseMapper {
	return &orderQueryResponseMapper{}
}

func (s *orderQueryResponseMapper) ToOrderResponse(order *record.OrderRecord) *response.OrderResponse {
	return &response.OrderResponse{
		ID:         order.ID,
		UserID:     order.UserID,
		TotalPrice: order.TotalPrice,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}
}

func (s *orderQueryResponseMapper) ToOrdersResponse(orders []*record.OrderRecord) []*response.OrderResponse {
	var responses []*response.OrderResponse

	for _, order := range orders {
		responses = append(responses, s.ToOrderResponse(order))
	}

	return responses
}

func (s *orderQueryResponseMapper) ToOrderResponseDeleteAt(order *record.OrderRecord) *response.OrderResponseDeleteAt {
	return &response.OrderResponseDeleteAt{
		ID:         order.ID,
		UserID:     order.UserID,
		TotalPrice: order.TotalPrice,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
		DeletedAt:  order.DeletedAt,
	}
}

func (s *orderQueryResponseMapper) ToOrdersResponseDeleteAt(orders []*record.OrderRecord) []*response.OrderResponseDeleteAt {
	var responses []*response.OrderResponseDeleteAt

	for _, order := range orders {
		responses = append(responses, s.ToOrderResponseDeleteAt(order))
	}

	return responses
}
