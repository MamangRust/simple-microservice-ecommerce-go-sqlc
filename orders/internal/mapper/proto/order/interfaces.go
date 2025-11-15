package orderprotomapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb"
	pborder "github.com/MamangRust/simple_microservice_ecommerce_pb/order"
)

type OrderBaseProtoMapper interface {
	ToProtoResponseOrder(status string, message string, pbResponse *response.OrderResponse) *pborder.ApiResponseOrder
}

type OrderQueryProtoMapper interface {
	ToProtoResponseOrder(status string, message string, pbResponse *response.OrderResponse) *pborder.ApiResponseOrder
	ToProtoResponseOrderDeleteAt(status string, message string, pbResponse *response.OrderResponseDeleteAt) *pborder.ApiResponseOrderDeleteAt

	ToProtoResponsePaginationOrderDeleteAt(pagination *pb.Pagination, status string, message string, Orders []*response.OrderResponseDeleteAt) *pborder.ApiResponsePaginationOrderDeleteAt
	ToProtoResponsePaginationOrder(pagination *pb.Pagination, status string, message string, Orders []*response.OrderResponse) *pborder.ApiResponsePaginationOrder
}

type OrderCommandProtoMapper interface {
	ToProtoResponseOrder(status string, message string, pbResponse *response.OrderResponse) *pborder.ApiResponseOrder
	ToProtoResponseOrderDeleteAt(status string, message string, pbResponse *response.OrderResponseDeleteAt) *pborder.ApiResponseOrderDeleteAt
	ToProtoResponseOrderDelete(status string, message string) *pborder.ApiResponseOrderDelete
	ToProtoResponseOrderAll(status string, message string) *pborder.ApiResponseOrderAll
}
