package repository

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/requests"
)

type ProductQueryRepository interface {
	FindAllProducts(ctx context.Context, req *requests.FindAllProduct) ([]*record.ProductRecord, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllProduct) ([]*record.ProductRecord, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*record.ProductRecord, *int, error)
	FindById(ctx context.Context, product_id int) (*record.ProductRecord, error)
}

type ProductCommandRepository interface {
	CreateProduct(ctx context.Context, request *requests.CreateProductRequest) (*record.ProductRecord, error)
	UpdateProduct(ctx context.Context, request *requests.UpdateProductRequest) (*record.ProductRecord, error)
	UpdateProductCountStock(ctx context.Context, req *requests.UpdateProductStockRequest) (*record.ProductRecord, error)
	TrashedProduct(ctx context.Context, product_id int) (*record.ProductRecord, error)
	RestoreProduct(ctx context.Context, product_id int) (*record.ProductRecord, error)
	DeleteProductPermanent(ctx context.Context, product_id int) (bool, error)
	RestoreAllProducts(ctx context.Context) (bool, error)
	DeleteAllProductPermanent(ctx context.Context) (bool, error)
}
