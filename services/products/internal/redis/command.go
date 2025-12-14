package mencache

import (
	"context"
	"fmt"
)

type productCommandCache struct {
	store *CacheStore
}

func NewProductCommandCache(store *CacheStore) *productCommandCache {
	return &productCommandCache{store: store}
}

func (c *productCommandCache) DeleteCachedProduct(ctx context.Context, productID int) {
	DeleteFromCache(ctx, c.store, fmt.Sprintf(productByIdCacheKey, productID))
}

func (s *productCommandCache) InvalidateAllProducts(ctx context.Context) {
	InvalidateByPattern(ctx, s.store, "product:all:*")
}

func (s *productCommandCache) InvalidateActiveProducts(ctx context.Context) {
	InvalidateByPattern(ctx, s.store, "product:active:*")
}

func (s *productCommandCache) InvalidateTrashedProducts(ctx context.Context) {
	InvalidateByPattern(ctx, s.store, "product:trashed:*")
}
