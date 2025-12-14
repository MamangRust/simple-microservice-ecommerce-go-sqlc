package mencache

import (
	"context"
	"fmt"

	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/domain/response"
)

const (
	productAllCacheKey     = "product:all:page:%d:pageSize:%d:search:%s"
	productByIdCacheKey    = "product:id:%d"
	productActiveCacheKey  = "product:active:page:%d:pageSize:%d:search:%s"
	productTrashedCacheKey = "product:trashed:page:%d:pageSize:%d:search:%s"
)

type productCache struct {
	store *CacheStore
}

func NewProductCache(store *CacheStore) *productCache {
	return &productCache{store: store}
}

func (s *productCache) GetCachedProducts(ctx context.Context, req *requests.FindAllProduct) (*response.ApiResponsePaginationProduct, bool) {
	key := fmt.Sprintf(productAllCacheKey, req.Page, req.PageSize, req.Search)

	payload, found := GetFromCache[response.ApiResponsePaginationProduct](ctx, s.store, key)

	if !found || payload.Data == nil {
		return nil, false
	}

	return payload, true
}

func (s *productCache) SetCachedProducts(
	ctx context.Context,
	req *requests.FindAllProduct,
	data *response.ApiResponsePaginationProduct,
) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(productAllCacheKey, req.Page, req.PageSize, req.Search)
	SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *productCache) GetCachedProductActive(ctx context.Context, req *requests.FindAllProduct) (*response.ApiResponsePaginationProductDeleteAt, bool) {
	key := fmt.Sprintf(productActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload, found := GetFromCache[response.ApiResponsePaginationProductDeleteAt](ctx, s.store, key)

	if !found || payload.Data == nil {
		return nil, false
	}

	return payload, true
}

func (s *productCache) SetCachedProductActive(ctx context.Context, req *requests.FindAllProduct, data *response.ApiResponsePaginationProductDeleteAt) {
	if data == nil {
		data = nil
	}

	key := fmt.Sprintf(productActiveCacheKey, req.Page, req.PageSize, req.Search)

	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *productCache) GetCachedProductTrashed(ctx context.Context, req *requests.FindAllProduct) (*response.ApiResponsePaginationProductDeleteAt, bool) {
	key := fmt.Sprintf(productTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload, found := GetFromCache[response.ApiResponsePaginationProductDeleteAt](ctx, s.store, key)

	if !found || payload.Data == nil {
		return nil, false
	}

	return payload, true
}

func (s *productCache) SetCachedProductTrashed(ctx context.Context, req *requests.FindAllProduct, data *response.ApiResponsePaginationProductDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(productTrashedCacheKey, req.Page, req.PageSize, req.Search)

	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *productCache) GetCachedProduct(ctx context.Context, productID int) (*response.ApiResponseProduct, bool) {
	key := fmt.Sprintf(productByIdCacheKey, productID)

	result, found := GetFromCache[response.ApiResponseProduct](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *productCache) SetCachedProduct(ctx context.Context, data *response.ApiResponseProduct) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(productByIdCacheKey, data.Data.ID)
	SetToCache(ctx, s.store, key, data, ttlDefault)
}
