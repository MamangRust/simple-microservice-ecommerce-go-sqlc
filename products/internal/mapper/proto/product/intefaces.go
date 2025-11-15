package productprotomapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/response"
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb"

	pbproduct "github.com/MamangRust/simple_microservice_ecommerce_pb/product"
)

type ProductBaseProtoMapper interface {
	ToProtoResponseProduct(status string, message string, pbResponse *response.ProductResponse) *pbproduct.ApiResponseProduct
}

type ProductQueryProtoMapper interface {
	ProductBaseProtoMapper

	ToProtoResponsePaginationProduct(pagination *pb.Pagination, status string, message string, pbResponse []*response.ProductResponse) *pbproduct.ApiResponsePaginationProduct

	ToProtoResponsePaginationProductDeleteAt(pagination *pb.Pagination, status string, message string, pbResponse []*response.ProductResponseDeleteAt) *pbproduct.ApiResponsePaginationProductDeleteAt
}

type ProductCommandProtoMapper interface {
	ProductBaseProtoMapper

	ToProtoResponseProductDeleteAt(status string, message string, pbResponse *response.ProductResponseDeleteAt) *pbproduct.ApiResponseProductDeleteAt

	ToProtoResponseProductAll(status string, message string) *pbproduct.ApiResponseProductAll

	ToProtoResponseProductDelete(status string, message string) *pbproduct.ApiResponseProductDelete
}
