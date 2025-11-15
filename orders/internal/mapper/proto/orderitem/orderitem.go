package orderitemprotomapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
	protomapper "github.com/MamangRust/simple_microservice_ecommerce/order/internal/mapper/proto"
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb"
	pborderitem "github.com/MamangRust/simple_microservice_ecommerce_pb/order_item"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type orderItemProtoMapper struct{}

func NewOrderItemProtoMapper() *orderItemProtoMapper {
	return &orderItemProtoMapper{}
}

func (o *orderItemProtoMapper) ToProtoResponseOrderItem(status string, message string, pbResponse *response.OrderItemResponse) *pborderitem.ApiResponseOrderItem {
	return &pborderitem.ApiResponseOrderItem{
		Status:  status,
		Message: message,
		Data:    o.mapResponseOrderItem(pbResponse),
	}
}

func (o *orderItemProtoMapper) ToProtoResponsesOrderItem(status string, message string, pbResponse []*response.OrderItemResponse) *pborderitem.ApiResponsesOrderItem {
	return &pborderitem.ApiResponsesOrderItem{
		Status:  status,
		Message: message,
		Data:    o.mapResponsesOrderItem(pbResponse),
	}
}

func (o *orderItemProtoMapper) ToProtoResponseOrderItemDelete(status string, message string) *pborderitem.ApiResponseOrderItemDelete {
	return &pborderitem.ApiResponseOrderItemDelete{
		Status:  status,
		Message: message,
	}
}

func (o *orderItemProtoMapper) ToProtoResponseOrderItemAll(status string, message string) *pborderitem.ApiResponseOrderItemAll {
	return &pborderitem.ApiResponseOrderItemAll{
		Status:  status,
		Message: message,
	}
}

func (o *orderItemProtoMapper) ToProtoResponsePaginationOrderItemDeleteAt(pagination *pb.Pagination, status string, message string, orderItems []*response.OrderItemResponseDeleteAt) *pborderitem.ApiResponsePaginationOrderItemDeleteAt {
	return &pborderitem.ApiResponsePaginationOrderItemDeleteAt{
		Status:     status,
		Message:    message,
		Data:       o.mapResponsesOrderItemDeleteAt(orderItems),
		Pagination: protomapper.MapPaginationMeta(pagination),
	}
}

func (o *orderItemProtoMapper) ToProtoResponsePaginationOrderItem(pagination *pb.Pagination, status string, message string, orderItems []*response.OrderItemResponse) *pborderitem.ApiResponsePaginationOrderItem {
	return &pborderitem.ApiResponsePaginationOrderItem{
		Status:     status,
		Message:    message,
		Data:       o.mapResponsesOrderItem(orderItems),
		Pagination: protomapper.MapPaginationMeta(pagination),
	}
}

func (o *orderItemProtoMapper) mapResponseOrderItem(orderItem *response.OrderItemResponse) *pborderitem.OrderItemResponse {
	return &pborderitem.OrderItemResponse{
		Id:        int32(orderItem.ID),
		OrderId:   int32(orderItem.OrderID),
		ProductId: int32(orderItem.ProductID),
		Quantity:  int32(orderItem.Quantity),
		Price:     int32(orderItem.Price),
		CreatedAt: orderItem.CreatedAt,
		UpdatedAt: orderItem.UpdatedAt,
	}
}

func (o *orderItemProtoMapper) mapResponsesOrderItem(orderItems []*response.OrderItemResponse) []*pborderitem.OrderItemResponse {
	var mappedOrderItems []*pborderitem.OrderItemResponse

	for _, orderItem := range orderItems {
		mappedOrderItems = append(mappedOrderItems, o.mapResponseOrderItem(orderItem))
	}

	return mappedOrderItems
}

func (o *orderItemProtoMapper) mapResponseOrderItemDelete(orderItem *response.OrderItemResponseDeleteAt) *pborderitem.OrderItemResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue

	if orderItem.DeletedAt != nil {
		deletedAt = wrapperspb.String(*orderItem.DeletedAt)
	}

	return &pborderitem.OrderItemResponseDeleteAt{
		Id:        int32(orderItem.ID),
		OrderId:   int32(orderItem.OrderID),
		ProductId: int32(orderItem.ProductID),
		Quantity:  int32(orderItem.Quantity),
		Price:     int32(orderItem.Price),
		CreatedAt: orderItem.CreatedAt,
		UpdatedAt: orderItem.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func (o *orderItemProtoMapper) mapResponsesOrderItemDeleteAt(orderItems []*response.OrderItemResponseDeleteAt) []*pborderitem.OrderItemResponseDeleteAt {
	var mappedOrderItems []*pborderitem.OrderItemResponseDeleteAt

	for _, orderItem := range orderItems {
		mappedOrderItems = append(mappedOrderItems, o.mapResponseOrderItemDelete(orderItem))
	}

	return mappedOrderItems
}
