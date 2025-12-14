package mencache

import (
	"context"
	"fmt"
)

type orderCommandCache struct {
	store *CacheStore
}

func NewOrderCommandCache(store *CacheStore) *orderCommandCache {
	return &orderCommandCache{store: store}
}

func (s *orderCommandCache) DeleteOrderCache(ctx context.Context, order_id int) {
	DeleteFromCache(ctx, s.store, fmt.Sprintf(orderByIdCacheKey, order_id))
}

func (s *orderCommandCache) InvalidateAllOrders(ctx context.Context) {
	InvalidateByPattern(ctx, s.store, "order:all:*")
}

func (s *orderCommandCache) InvalidateActiveOrders(ctx context.Context) {
	InvalidateByPattern(ctx, s.store, "order:active:*")
}

func (s *orderCommandCache) InvalidateTrashedOrders(ctx context.Context) {
	InvalidateByPattern(ctx, s.store, "order:trashed:*")
}
