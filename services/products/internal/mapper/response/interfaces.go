package productresponsemapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/response"
)

type ProductBaseResponseMapper interface {
	ToProductResponse(Product *record.ProductRecord) *response.ProductResponse
}

type ProductQueryResponseMapper interface {
	ProductBaseResponseMapper

	ToProductsResponse(Products []*record.ProductRecord) []*response.ProductResponse
	ToProductsResponseDeleteAt(Products []*record.ProductRecord) []*response.ProductResponseDeleteAt
}

type ProductCommandResponseMapper interface {
	ProductBaseResponseMapper

	ToProductResponseDeleteAt(Product *record.ProductRecord) *response.ProductResponseDeleteAt
}
