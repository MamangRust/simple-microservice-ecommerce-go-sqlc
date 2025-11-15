package productresponsemapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/response"
)

type productCommandResponseMapper struct{}

func NewProductCommandResponseMapper() ProductCommandResponseMapper {
	return &productCommandResponseMapper{}
}

func (s *productCommandResponseMapper) ToProductResponse(product *record.ProductRecord) *response.ProductResponse {
	return &response.ProductResponse{
		ID:        product.ID,
		Name:      product.Name,
		Price:     product.Price,
		Stock:     product.Stock,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}

func (s *productCommandResponseMapper) ToProductResponseDeleteAt(product *record.ProductRecord) *response.ProductResponseDeleteAt {
	return &response.ProductResponseDeleteAt{
		ID:        product.ID,
		Name:      product.Name,
		Price:     product.Price,
		Stock:     product.Stock,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
		DeletedAt: product.DeletedAt,
	}
}
