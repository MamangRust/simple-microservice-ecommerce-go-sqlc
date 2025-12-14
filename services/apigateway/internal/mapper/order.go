package mapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/domain/response"
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb/order"
)

type orderResponseMapper struct {
}

func NewOrderResponseMapper() OrderResponseMapper {
	return &orderResponseMapper{}
}

func (o *orderResponseMapper) ToResponseOrder(order *pb.OrderResponse) *response.OrderResponse {
	return &response.OrderResponse{
		ID:         int(order.Id),
		UserID:     int(order.UserId),
		TotalPrice: int(order.TotalPrice),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}
}

func (o *orderResponseMapper) ToResponsesOrder(orders []*pb.OrderResponse) []*response.OrderResponse {
	var mappedOrders []*response.OrderResponse

	for _, order := range orders {
		mappedOrders = append(mappedOrders, o.ToResponseOrder(order))
	}

	return mappedOrders
}

func (o *orderResponseMapper) ToResponseOrderDeleteAt(order *pb.OrderResponseDeleteAt) *response.OrderResponseDeleteAt {
	var deletedAt string
	if order.DeletedAt != nil {
		deletedAt = order.DeletedAt.Value
	}

	return &response.OrderResponseDeleteAt{
		ID:         int(order.Id),
		UserID:     int(order.UserId),
		TotalPrice: int(order.TotalPrice),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
		DeletedAt:  &deletedAt,
	}
}

func (o *orderResponseMapper) ToResponsesOrderDeleteAt(orders []*pb.OrderResponseDeleteAt) []*response.OrderResponseDeleteAt {
	var mappedOrders []*response.OrderResponseDeleteAt

	for _, order := range orders {
		mappedOrders = append(mappedOrders, o.ToResponseOrderDeleteAt(order))
	}

	return mappedOrders
}

func (o *orderResponseMapper) ToApiResponseOrder(pbResponse *pb.ApiResponseOrder) *response.ApiResponseOrder {
	return &response.ApiResponseOrder{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    o.ToResponseOrder(pbResponse.Data),
	}
}

func (o *orderResponseMapper) ToApiResponseOrderDeleteAt(pbResponse *pb.ApiResponseOrderDeleteAt) *response.ApiResponseOrderDeleteAt {
	return &response.ApiResponseOrderDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    o.ToResponseOrderDeleteAt(pbResponse.Data),
	}
}

func (o *orderResponseMapper) ToApiResponseOrderDelete(pbResponse *pb.ApiResponseOrderDelete) *response.ApiResponseOrderDelete {
	return &response.ApiResponseOrderDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (o *orderResponseMapper) ToApiResponseOrderAll(pbResponse *pb.ApiResponseOrderAll) *response.ApiResponseOrderAll {
	return &response.ApiResponseOrderAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (o *orderResponseMapper) ToApiResponsePaginationOrderDeleteAt(pbResponse *pb.ApiResponsePaginationOrderDeleteAt) *response.ApiResponsePaginationOrderDeleteAt {
	return &response.ApiResponsePaginationOrderDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       o.ToResponsesOrderDeleteAt(pbResponse.Data),
		Pagination: mapPaginationMeta(pbResponse.Pagination),
	}
}

func (o *orderResponseMapper) ToApiResponsePaginationOrder(pbResponse *pb.ApiResponsePaginationOrder) *response.ApiResponsePaginationOrder {
	return &response.ApiResponsePaginationOrder{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       o.ToResponsesOrder(pbResponse.Data),
		Pagination: mapPaginationMeta(pbResponse.Pagination),
	}
}
