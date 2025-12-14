package mencache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/domain/response"
)

const (
	orderAllCacheKey     = "order:all:page:%d:pageSize:%d:search:%s"
	orderByIdCacheKey    = "order:id:%d"
	orderActiveCacheKey  = "order:active:page:%d:pageSize:%d:search:%s"
	orderTrashedCacheKey = "order:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type orderCache struct {
	store *CacheStore
}

func NewOrderCache(store *CacheStore) *orderCache {
	return &orderCache{store: store}
}

func (s *orderCache) GetOrderAllCache(ctx context.Context, req *requests.FindAllOrder) (*response.ApiResponsePaginationOrder, bool) {
	key := fmt.Sprintf(orderAllCacheKey, req.Page, req.PageSize, req.Search)

	payload, found := GetFromCache[response.ApiResponsePaginationOrder](ctx, s.store, key)

	if !found || payload.Data == nil {
		return nil, false
	}

	return payload, true
}

func (s *orderCache) SetOrderAllCache(ctx context.Context, req *requests.FindAllOrder, data *response.ApiResponsePaginationOrder) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(orderAllCacheKey, req.Page, req.PageSize, req.Search)

	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *orderCache) GetOrderActiveCache(ctx context.Context, req *requests.FindAllOrder) (*response.ApiResponsePaginationOrderDeleteAt, bool) {
	key := fmt.Sprintf(orderActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload, found := GetFromCache[response.ApiResponsePaginationOrderDeleteAt](ctx, s.store, key)

	if !found || payload.Data == nil {
		return nil, false
	}

	return payload, true
}

func (s *orderCache) SetOrderActiveCache(ctx context.Context, req *requests.FindAllOrder, data *response.ApiResponsePaginationOrderDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(orderActiveCacheKey, req.Page, req.PageSize, req.Search)

	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *orderCache) GetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrder) (*response.ApiResponsePaginationOrderDeleteAt, bool) {
	key := fmt.Sprintf(orderTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload, found := GetFromCache[response.ApiResponsePaginationOrderDeleteAt](ctx, s.store, key)

	if !found || payload.Data == nil {
		return nil, false
	}

	return payload, true
}

func (s *orderCache) SetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrder, data *response.ApiResponsePaginationOrderDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(orderTrashedCacheKey, req.Page, req.PageSize, req.Search)

	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *orderCache) GetCachedOrderCache(ctx context.Context, orderID int) (*response.ApiResponseOrder, bool) {
	key := fmt.Sprintf(orderByIdCacheKey, orderID)

	result, found := GetFromCache[response.ApiResponseOrder](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *orderCache) SetCachedOrderCache(ctx context.Context, data *response.ApiResponseOrder) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(orderByIdCacheKey, data.Data.ID)
	SetToCache(ctx, s.store, key, data, ttlDefault)
}
