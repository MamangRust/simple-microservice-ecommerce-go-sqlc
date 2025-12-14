package mencache

import (
	"context"
)

type userCommandCache struct {
	store *CacheStore
}

func NewUserCommandCache(store *CacheStore) UserCommandCache {
	return &userCommandCache{store: store}
}

func (s *userCommandCache) DeleteCachedUser(ctx context.Context, key string) {
	DeleteFromCache(ctx, s.store, key)
}

func (s *userCommandCache) InvalidateAllUsers(ctx context.Context) {
	InvalidateByPattern(ctx, s.store, "user:all:*")
}

func (s *userCommandCache) InvalidateActiveUsers(ctx context.Context) {
	InvalidateByPattern(ctx, s.store, "user:active:*")
}

func (s *userCommandCache) InvalidateTrashedUsers(ctx context.Context) {
	InvalidateByPattern(ctx, s.store, "user:trashed:*")
}
