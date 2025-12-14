package mencache

import (
	"context"
	"fmt"
)

type roleCommandCache struct {
	store *CacheStore
}

func NewRoleCommandCache(store *CacheStore) *roleCommandCache {
	return &roleCommandCache{store: store}
}

func (s *roleCommandCache) DeleteCachedRole(ctx context.Context, id int) {
	key := fmt.Sprintf(roleByIdCacheKey, id)

	DeleteFromCache(ctx, s.store, key)
}

func (s *roleCommandCache) InvalidateAllRoles(ctx context.Context) {
	InvalidateByPattern(ctx, s.store, "role:all:*")
}

func (s *roleCommandCache) InvalidateActiveRoles(ctx context.Context) {
	InvalidateByPattern(ctx, s.store, "role:active:*")
}

func (s *roleCommandCache) InvalidateTrashedRoles(ctx context.Context) {
	InvalidateByPattern(ctx, s.store, "role:trashed:*")
}
