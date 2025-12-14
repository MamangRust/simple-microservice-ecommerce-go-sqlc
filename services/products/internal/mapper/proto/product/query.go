package productprotomapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/response"
	helperproto "github.com/MamangRust/simple_microservice_ecommerce/product/internal/mapper/proto/helpers"
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb"
	pbproduct "github.com/MamangRust/simple_microservice_ecommerce_pb/product"
)

type productQueryProtoMapper struct {
}

func NewProductQueryProtoMapper() ProductQueryProtoMapper {
	return &productQueryProtoMapper{}
}

func (p *productQueryProtoMapper) ToProtoResponseProduct(status string, message string, pbResponse *response.ProductResponse) *pbproduct.ApiResponseProduct {
	return &pbproduct.ApiResponseProduct{
		Status:  status,
		Message: message,
		Data:    p.mapResponseProduct(pbResponse),
	}
}

func (p *productQueryProtoMapper) ToProtoResponseProductDeleteAt(status string, message string, pbResponse *response.ProductResponseDeleteAt) *pbproduct.ApiResponseProductDeleteAt {
	return &pbproduct.ApiResponseProductDeleteAt{
		Status:  status,
		Message: message,
		Data:    p.mapResponseProductDeleteAt(pbResponse),
	}
}

func (p *productQueryProtoMapper) ToProtoResponsePaginationProductDeleteAt(pagination *pb.Pagination, status string, message string, products []*response.ProductResponseDeleteAt) *pbproduct.ApiResponsePaginationProductDeleteAt {
	return &pbproduct.ApiResponsePaginationProductDeleteAt{
		Status:     status,
		Message:    message,
		Data:       p.mapResponsesProductDeleteAt(products),
		Pagination: MapPaginationMeta(pagination),
	}
}

func (p *productQueryProtoMapper) ToProtoResponsePaginationProduct(pagination *pb.Pagination, status string, message string, products []*response.ProductResponse) *pbproduct.ApiResponsePaginationProduct {
	return &pbproduct.ApiResponsePaginationProduct{
		Status:     status,
		Message:    message,
		Data:       p.mapResponsesProduct(products),
		Pagination: MapPaginationMeta(pagination),
	}
}

func (p *productQueryProtoMapper) mapResponseProduct(product *response.ProductResponse) *pbproduct.ProductResponse {
	return &pbproduct.ProductResponse{
		Id:        int32(product.ID),
		Name:      product.Name,
		Price:     int64(product.Price),
		Stock:     int32(product.Stock),
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}

func (p *productQueryProtoMapper) mapResponsesProduct(products []*response.ProductResponse) []*pbproduct.ProductResponse {
	var mappedProducts []*pbproduct.ProductResponse

	for _, product := range products {
		mappedProducts = append(mappedProducts, p.mapResponseProduct(product))
	}

	return mappedProducts
}

func (p *productQueryProtoMapper) mapResponseProductDeleteAt(product *response.ProductResponseDeleteAt) *pbproduct.ProductResponseDeleteAt {
	res := &pbproduct.ProductResponseDeleteAt{
		Id:        int32(product.ID),
		Name:      product.Name,
		Price:     int64(product.Price),
		Stock:     int32(product.Stock),
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}

	if product.DeletedAt != nil {
		res.DeletedAt = helperproto.StringPtrToWrapper(product.DeletedAt)
	}

	return res
}

func (p *productQueryProtoMapper) mapResponsesProductDeleteAt(
	products []*response.ProductResponseDeleteAt,
) []*pbproduct.ProductResponseDeleteAt {
	result := make([]*pbproduct.ProductResponseDeleteAt, 0, len(products))
	for _, prod := range products {
		result = append(result, p.mapResponseProductDeleteAt(prod))
	}
	return result
}

func MapPaginationMeta(s *pb.Pagination) *pb.Pagination {
	return &pb.Pagination{
		CurrentPage:  int32(s.CurrentPage),
		PageSize:     int32(s.PageSize),
		TotalPages:   int32(s.TotalPages),
		TotalRecords: int32(s.TotalRecords),
	}
}
