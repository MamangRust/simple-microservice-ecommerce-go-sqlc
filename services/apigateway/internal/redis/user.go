package mencache

import (
	"context"
	"fmt"

	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/domain/response"
)

const (
	userAllCacheKey     = "user:all:page:%d:pageSize:%d:search:%s"
	userByIdCacheKey    = "user:id:%d"
	userActiveCacheKey  = "user:active:page:%d:pageSize:%d:search:%s"
	userTrashedCacheKey = "user:trashed:page:%d:pageSize:%d:search:%s"
)

type userCache struct {
	store *CacheStore
}

func NewUserCache(store *CacheStore) *userCache {
	return &userCache{store: store}
}

func (s *userCache) SetCachedUsers(ctx context.Context, req *requests.FindAllUsers, data *response.ApiResponsePaginationUser) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(userAllCacheKey, req.Page, req.PageSize, req.Search)

	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *userCache) SetCachedUserById(ctx context.Context, data *response.ApiResponseUser) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(userByIdCacheKey, data.Data.ID)
	SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *userCache) SetCachedUserActive(ctx context.Context, req *requests.FindAllUsers, data *response.ApiResponsePaginationUserDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(userActiveCacheKey, req.Page, req.PageSize, req.Search)

	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *userCache) SetCachedUserTrashed(ctx context.Context, req *requests.FindAllUsers, data *response.ApiResponsePaginationUserDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(userTrashedCacheKey, req.Page, req.PageSize, req.Search)

	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *userCache) GetCachedUsers(ctx context.Context, req *requests.FindAllUsers) (*response.ApiResponsePaginationUser, bool) {
	key := fmt.Sprintf(userAllCacheKey, req.Page, req.PageSize, req.Search)

	payload, found := GetFromCache[response.ApiResponsePaginationUser](ctx, s.store, key)

	if !found || payload.Data == nil {
		return nil, false
	}

	return payload, true
}

func (s *userCache) GetCachedUserById(ctx context.Context, id int) (*response.ApiResponseUser, bool) {
	key := fmt.Sprintf(userByIdCacheKey, id)

	result, found := GetFromCache[response.ApiResponseUser](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *userCache) GetCachedUserActive(ctx context.Context, req *requests.FindAllUsers) (*response.ApiResponsePaginationUserDeleteAt, bool) {
	key := fmt.Sprintf(userActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload, found := GetFromCache[response.ApiResponsePaginationUserDeleteAt](ctx, s.store, key)

	if !found || payload.Data == nil {
		return nil, false
	}

	return payload, true
}

func (s *userCache) GetCachedUserTrashed(ctx context.Context, req *requests.FindAllUsers) (*response.ApiResponsePaginationUserDeleteAt, bool) {
	key := fmt.Sprintf(userTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload, found := GetFromCache[response.ApiResponsePaginationUserDeleteAt](ctx, s.store, key)

	if !found || payload.Data == nil {
		return nil, false
	}

	return payload, true
}
