package productresponsemapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/response"
)

type productQueryResponseMapper struct{}

func NewProductQueryResponseMapper() ProductQueryResponseMapper {
	return &productQueryResponseMapper{}
}

func (s *productQueryResponseMapper) ToProductResponse(product *record.ProductRecord) *response.ProductResponse {
	return &response.ProductResponse{
		ID:        product.ID,
		Name:      product.Name,
		Price:     product.Price,
		Stock:     product.Stock,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}

func (s *productQueryResponseMapper) ToProductsResponse(products []*record.ProductRecord) []*response.ProductResponse {
	var responseProducts []*response.ProductResponse

	for _, product := range products {
		responseProducts = append(responseProducts, s.ToProductResponse(product))
	}

	return responseProducts
}

func (s *productQueryResponseMapper) ToProductResponseDeleteAt(product *record.ProductRecord) *response.ProductResponseDeleteAt {
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

func (s *productQueryResponseMapper) ToProductsResponseDeleteAt(products []*record.ProductRecord) []*response.ProductResponseDeleteAt {
	var responseProducts []*response.ProductResponseDeleteAt

	for _, product := range products {
		responseProducts = append(responseProducts, s.ToProductResponseDeleteAt(product))
	}

	return responseProducts
}
