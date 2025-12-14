package mapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/domain/response"
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb/product"
)

type productResponseMapper struct {
}

func NewProductResponseMapper() ProductResponseMapper {
	return &productResponseMapper{}
}

func (p *productResponseMapper) ToResponseProduct(product *pb.ProductResponse) *response.ProductResponse {
	return &response.ProductResponse{
		ID:        int(product.Id),
		Name:      product.Name,
		Price:     int(product.Price),
		Stock:     int(product.Stock),
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}

func (p *productResponseMapper) ToResponsesProduct(products []*pb.ProductResponse) []*response.ProductResponse {
	var mappedProducts []*response.ProductResponse

	for _, product := range products {
		mappedProducts = append(mappedProducts, p.ToResponseProduct(product))
	}

	return mappedProducts
}

func (p *productResponseMapper) ToResponseProductDeleteAt(product *pb.ProductResponseDeleteAt) *response.ProductResponseDeleteAt {
	var deletedAt string
	if product.DeletedAt != nil {
		deletedAt = product.DeletedAt.Value
	}

	return &response.ProductResponseDeleteAt{
		ID:        int(product.Id),
		Name:      product.Name,
		Price:     int(product.Price),
		Stock:     int(product.Stock),
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
		DeletedAt: &deletedAt,
	}
}

func (p *productResponseMapper) ToResponsesProductDeleteAt(products []*pb.ProductResponseDeleteAt) []*response.ProductResponseDeleteAt {
	var mappedProducts []*response.ProductResponseDeleteAt

	for _, product := range products {
		mappedProducts = append(mappedProducts, p.ToResponseProductDeleteAt(product))
	}

	return mappedProducts
}

func (p *productResponseMapper) ToApiResponseProduct(pbResponse *pb.ApiResponseProduct) *response.ApiResponseProduct {
	return &response.ApiResponseProduct{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    p.ToResponseProduct(pbResponse.Data),
	}
}

func (p *productResponseMapper) ToApiResponsesProductDeleteAt(pbResponse *pb.ApiResponseProductDeleteAt) *response.ApiResponseProductDeleteAt {
	return &response.ApiResponseProductDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    p.ToResponseProductDeleteAt(pbResponse.Data),
	}
}

func (p *productResponseMapper) ToApiResponseProductDelete(pbResponse *pb.ApiResponseProductDelete) *response.ApiResponseProductDelete {
	return &response.ApiResponseProductDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (p *productResponseMapper) ToApiResponseProductAll(pbResponse *pb.ApiResponseProductAll) *response.ApiResponseProductAll {
	return &response.ApiResponseProductAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (p *productResponseMapper) ToApiResponsePaginationProductDeleteAt(pbResponse *pb.ApiResponsePaginationProductDeleteAt) *response.ApiResponsePaginationProductDeleteAt {
	return &response.ApiResponsePaginationProductDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       p.ToResponsesProductDeleteAt(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}

func (p *productResponseMapper) ToApiResponsePaginationProduct(pbResponse *pb.ApiResponsePaginationProduct) *response.ApiResponsePaginationProduct {
	return &response.ApiResponsePaginationProduct{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       p.ToResponsesProduct(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}
