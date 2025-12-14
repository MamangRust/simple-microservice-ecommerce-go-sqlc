package mencache

import (
	"context"
	"fmt"

	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/domain/response"
)

const (
	orderItemAllCacheKey       = "orderitem:all:page:%d:pageSize:%d:search:%s"
	orderItemActiveCacheKey    = "orderitem:active:page:%d:pageSize:%d:search:%s"
	orderItemTrashedCacheKey   = "orderitem:trashed:page:%d:pageSize:%d:search:%s"
	orderItemByOrderIdCacheKey = "orderitem:byorderid:%d"
)

type orderItemCache struct {
	store *CacheStore
}

func NewOrderItemCache(store *CacheStore) *orderItemCache {
	return &orderItemCache{store: store}
}

func (s *orderItemCache) GetCachedOrderItemsAll(ctx context.Context, req *requests.FindAllOrderItems) (*response.ApiResponsePaginationOrderItem, bool) {
	key := fmt.Sprintf(orderItemAllCacheKey, req.Page, req.PageSize, req.Search)

	payload, found := GetFromCache[response.ApiResponsePaginationOrderItem](ctx, s.store, key)

	if !found || payload.Data == nil {
		return nil, false
	}

	return payload, true
}

func (s *orderItemCache) SetCachedOrderItemsAll(ctx context.Context, req *requests.FindAllOrderItems, data *response.ApiResponsePaginationOrderItem) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(orderItemAllCacheKey, req.Page, req.PageSize, req.Search)
	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *orderItemCache) GetCachedOrderItemActive(ctx context.Context, req *requests.FindAllOrderItems) (*response.ApiResponsePaginationOrderItemDeleteAt, bool) {
	key := fmt.Sprintf(orderItemActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload, found := GetFromCache[response.ApiResponsePaginationOrderItemDeleteAt](ctx, s.store, key)

	if !found || payload.Data == nil {
		return nil, false
	}

	return payload, true
}

func (s *orderItemCache) SetCachedOrderItemActive(ctx context.Context, req *requests.FindAllOrderItems, data *response.ApiResponsePaginationOrderItemDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(orderItemActiveCacheKey, req.Page, req.PageSize, req.Search)

	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *orderItemCache) GetCachedOrderItemTrashed(ctx context.Context, req *requests.FindAllOrderItems) (*response.ApiResponsePaginationOrderItemDeleteAt, bool) {
	key := fmt.Sprintf(orderItemTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload, found := GetFromCache[response.ApiResponsePaginationOrderItemDeleteAt](ctx, s.store, key)

	if !found || payload.Data == nil {
		return nil, false
	}

	return payload, true
}

func (s *orderItemCache) SetCachedOrderItemTrashed(ctx context.Context, req *requests.FindAllOrderItems, data *response.ApiResponsePaginationOrderItemDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(orderItemTrashedCacheKey, req.Page, req.PageSize, req.Search)

	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *orderItemCache) GetCachedOrderItems(ctx context.Context, orderID int) (*response.ApiResponsesOrderItem, bool) {
	key := fmt.Sprintf(orderItemByOrderIdCacheKey, orderID)

	result, found := GetFromCache[*response.ApiResponsesOrderItem](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderItemCache) SetCachedOrderItems(ctx context.Context, orderID int, data *response.ApiResponsesOrderItem) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(orderItemByOrderIdCacheKey, orderID)

	SetToCache(ctx, s.store, key, &data, ttlDefault)
}
