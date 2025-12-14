package mencache

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/response"
)

type ProductQueryCache interface {
	GetCachedProducts(ctx context.Context, req *requests.FindAllProduct) ([]*response.ProductResponse, *int, bool)
	SetCachedProducts(ctx context.Context, req *requests.FindAllProduct, data []*response.ProductResponse, total *int)

	GetCachedProductActive(ctx context.Context, req *requests.FindAllProduct) ([]*response.ProductResponseDeleteAt, *int, bool)
	SetCachedProductActive(ctx context.Context, req *requests.FindAllProduct, data []*response.ProductResponseDeleteAt, total *int)

	GetCachedProductTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*response.ProductResponseDeleteAt, *int, bool)
	SetCachedProductTrashed(ctx context.Context, req *requests.FindAllProduct, data []*response.ProductResponseDeleteAt, total *int)

	GetCachedProduct(ctx context.Context, productID int) (*response.ProductResponse, bool)
	SetCachedProduct(ctx context.Context, data *response.ProductResponse)
}

type ProductCommandCache interface {
	DeleteCachedProduct(ctx context.Context, productID int)
	InvalidateAllProducts(ctx context.Context)
	InvalidateActiveProducts(ctx context.Context)
	InvalidateTrashedProducts(ctx context.Context)
}
