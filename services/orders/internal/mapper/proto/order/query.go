package orderprotomapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
	protomapper "github.com/MamangRust/simple_microservice_ecommerce/order/internal/mapper/proto"
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb"
	pborder "github.com/MamangRust/simple_microservice_ecommerce_pb/order"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type orderQueryProtoMapper struct {
}

func NewOrderQueryProtoMapper() OrderQueryProtoMapper {
	return &orderQueryProtoMapper{}
}

func (u *orderQueryProtoMapper) ToProtoResponseOrder(status string, message string, pbResponse *response.OrderResponse) *pborder.ApiResponseOrder {
	return &pborder.ApiResponseOrder{
		Status:  status,
		Message: message,
		Data:    u.mapResponseOrder(pbResponse),
	}
}

func (u *orderQueryProtoMapper) ToProtoResponseOrderDeleteAt(status string, message string, pbResponse *response.OrderResponseDeleteAt) *pborder.ApiResponseOrderDeleteAt {
	return &pborder.ApiResponseOrderDeleteAt{
		Status:  status,
		Message: message,
		Data:    u.mapResponseOrderDeleteAt(pbResponse),
	}
}

func (u *orderQueryProtoMapper) ToProtoResponsePaginationOrder(pagination *pb.Pagination, status string, message string, users []*response.OrderResponse) *pborder.ApiResponsePaginationOrder {
	return &pborder.ApiResponsePaginationOrder{
		Status:     status,
		Message:    message,
		Data:       u.mapResponsesOrder(users),
		Pagination: protomapper.MapPaginationMeta(pagination),
	}
}

func (u *orderQueryProtoMapper) ToProtoResponsePaginationOrderDeleteAt(pagination *pb.Pagination, status string, message string, users []*response.OrderResponseDeleteAt) *pborder.ApiResponsePaginationOrderDeleteAt {
	return &pborder.ApiResponsePaginationOrderDeleteAt{
		Status:     status,
		Message:    message,
		Data:       u.mapResponsesOrderDeleteAt(users),
		Pagination: protomapper.MapPaginationMeta(pagination),
	}
}

func (o *orderQueryProtoMapper) mapResponseOrder(order *response.OrderResponse) *pborder.OrderResponse {
	return &pborder.OrderResponse{
		Id:         int32(order.ID),
		UserId:     int32(order.UserID),
		TotalPrice: int32(order.TotalPrice),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}
}

func (o *orderQueryProtoMapper) mapResponsesOrder(orders []*response.OrderResponse) []*pborder.OrderResponse {
	var mappedOrders []*pborder.OrderResponse

	for _, order := range orders {
		mappedOrders = append(mappedOrders, o.mapResponseOrder(order))
	}

	return mappedOrders
}

func (o *orderQueryProtoMapper) mapResponseOrderDeleteAt(order *response.OrderResponseDeleteAt) *pborder.OrderResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue

	if order.DeletedAt != nil {
		deletedAt = wrapperspb.String(*order.DeletedAt)
	}

	return &pborder.OrderResponseDeleteAt{
		Id:         int32(order.ID),
		UserId:     int32(order.UserID),
		TotalPrice: int32(order.TotalPrice),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
		DeletedAt:  deletedAt,
	}
}

func (o *orderQueryProtoMapper) mapResponsesOrderDeleteAt(orders []*response.OrderResponseDeleteAt) []*pborder.OrderResponseDeleteAt {
	var mappedOrders []*pborder.OrderResponseDeleteAt

	for _, order := range orders {
		mappedOrders = append(mappedOrders, o.mapResponseOrderDeleteAt(order))
	}

	return mappedOrders
}
