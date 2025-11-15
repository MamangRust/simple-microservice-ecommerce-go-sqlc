package orderprotomapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
	pborder "github.com/MamangRust/simple_microservice_ecommerce_pb/order"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type orderCommandProtoMapper struct {
}

func NewOrderCommandProtoMapper() OrderCommandProtoMapper {
	return &orderCommandProtoMapper{}
}

func (u *orderCommandProtoMapper) ToProtoResponseOrder(status string, message string, pbResponse *response.OrderResponse) *pborder.ApiResponseOrder {
	return &pborder.ApiResponseOrder{
		Status:  status,
		Message: message,
		Data:    u.mapResponseOrder(pbResponse),
	}
}

func (u *orderCommandProtoMapper) ToProtoResponseOrderDeleteAt(status string, message string, pbResponse *response.OrderResponseDeleteAt) *pborder.ApiResponseOrderDeleteAt {
	return &pborder.ApiResponseOrderDeleteAt{
		Status:  status,
		Message: message,
		Data:    u.mapResponseOrderDeleteAt(pbResponse),
	}
}

func (o *orderCommandProtoMapper) ToProtoResponseOrderDelete(status string, message string) *pborder.ApiResponseOrderDelete {
	return &pborder.ApiResponseOrderDelete{
		Status:  status,
		Message: message,
	}
}

func (o *orderCommandProtoMapper) ToProtoResponseOrderAll(status string, message string) *pborder.ApiResponseOrderAll {
	return &pborder.ApiResponseOrderAll{
		Status:  status,
		Message: message,
	}
}

func (o *orderCommandProtoMapper) mapResponseOrder(order *response.OrderResponse) *pborder.OrderResponse {
	return &pborder.OrderResponse{
		Id:         int32(order.ID),
		UserId:     int32(order.UserID),
		TotalPrice: int32(order.TotalPrice),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}
}

func (o *orderCommandProtoMapper) mapResponseOrderDeleteAt(order *response.OrderResponseDeleteAt) *pborder.OrderResponseDeleteAt {
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
