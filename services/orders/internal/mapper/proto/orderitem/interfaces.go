package orderitemprotomapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb"
	pborderitem "github.com/MamangRust/simple_microservice_ecommerce_pb/order_item"
)

type OrderItemProtoMapper interface {
	ToProtoResponseOrderItem(status string, message string, pbResponse *response.OrderItemResponse) *pborderitem.ApiResponseOrderItem
	ToProtoResponsesOrderItem(status string, message string, pbResponse []*response.OrderItemResponse) *pborderitem.ApiResponsesOrderItem
	ToProtoResponseOrderItemDelete(status string, message string) *pborderitem.ApiResponseOrderItemDelete
	ToProtoResponseOrderItemAll(status string, message string) *pborderitem.ApiResponseOrderItemAll
	ToProtoResponsePaginationOrderItemDeleteAt(pagination *pb.Pagination, status string, message string, orderItems []*response.OrderItemResponseDeleteAt) *pborderitem.ApiResponsePaginationOrderItemDeleteAt
	ToProtoResponsePaginationOrderItem(pagination *pb.Pagination, status string, message string, orderItems []*response.OrderItemResponse) *pborderitem.ApiResponsePaginationOrderItem
}
