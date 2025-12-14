package grpcclientmapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
	pbproduct "github.com/MamangRust/simple_microservice_ecommerce_pb/product"
)

type ProductClientResponseMapper interface {
	ToApiResponseProduct(pbResponse *pbproduct.ApiResponseProduct) *response.ApiResponseProduct
}

type productClientResponseMapper struct{}

func NewProductClientResponseMapper() ProductClientResponseMapper {
	return &productClientResponseMapper{}
}

func (p *productClientResponseMapper) ToResponseProduct(product *pbproduct.ProductResponse) *response.ProductResponse {
	return &response.ProductResponse{
		ID:        int(product.Id),
		Name:      product.Name,
		Price:     int(product.Price),
		Stock:     int(product.Stock),
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}

func (p *productClientResponseMapper) ToApiResponseProduct(pbResponse *pbproduct.ApiResponseProduct) *response.ApiResponseProduct {
	return &response.ApiResponseProduct{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    p.ToResponseProduct(pbResponse.Data),
	}
}
