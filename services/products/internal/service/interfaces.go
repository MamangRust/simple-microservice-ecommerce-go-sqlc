package service

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/response"
)

type ProductQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllProduct) ([]*response.ProductResponse, *int, *response.ErrorResponse)
	FindByID(ctx context.Context, id int) (*response.ProductResponse, *response.ErrorResponse)
	FindByActive(ctx context.Context, req *requests.FindAllProduct) ([]*response.ProductResponseDeleteAt, *int, *response.ErrorResponse)

	FindByTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*response.ProductResponseDeleteAt, *int, *response.ErrorResponse)
}

type ProductCommandService interface {
	CreateProduct(ctx context.Context, request *requests.CreateProductRequest) (*response.ProductResponse, *response.ErrorResponse)
	UpdateProduct(ctx context.Context, request *requests.UpdateProductRequest) (*response.ProductResponse, *response.ErrorResponse)
	UpdateProductStock(ctx context.Context, request *requests.UpdateProductStockRequest) (*response.ProductResponse, *response.ErrorResponse)
	TrashedProduct(ctx context.Context, Product_id int) (*response.ProductResponseDeleteAt, *response.ErrorResponse)
	RestoreProduct(ctx context.Context, Product_id int) (*response.ProductResponseDeleteAt, *response.ErrorResponse)
	DeleteProductPermanent(ctx context.Context, Product_id int) (bool, *response.ErrorResponse)
	RestoreAllProduct(ctx context.Context) (bool, *response.ErrorResponse)
	DeleteAllProductPermanent(ctx context.Context) (bool, *response.ErrorResponse)
}
