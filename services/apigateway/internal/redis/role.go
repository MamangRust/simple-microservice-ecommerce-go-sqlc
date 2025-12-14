package mencache

import (
	"context"
	"fmt"

	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/domain/response"
)

const (
	roleAllCacheKey      = "role:all:page:%d:pageSize:%d:search:%s"
	roleByIdCacheKey     = "role:id:%d"
	roleByUserIdCacheKey = "role:byuserid:%d"
	roleActiveCacheKey   = "role:active:page:%d:pageSize:%d:search:%s"
	roleTrashedCacheKey  = "role:trashed:page:%d:pageSize:%d:search:%s"
)

type roleCache struct {
	store *CacheStore
}

func NewRoleCache(store *CacheStore) *roleCache {
	return &roleCache{store: store}
}

func (s *roleCache) SetCachedRoles(ctx context.Context, req *requests.FindAllRoles, data *response.ApiResponsePaginationRole) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(roleAllCacheKey, req.Page, req.PageSize, req.Search)

	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *roleCache) SetCachedRoleById(ctx context.Context, data *response.ApiResponseRole) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(roleByIdCacheKey, data.Data.ID)
	SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *roleCache) SetCachedRoleByUserId(ctx context.Context, userId int, data *response.ApiResponsesRole) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(roleByUserIdCacheKey, userId)
	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *roleCache) SetCachedRoleActive(ctx context.Context, req *requests.FindAllRoles, data *response.ApiResponsePaginationRoleDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(roleActiveCacheKey, req.Page, req.PageSize, req.Search)

	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *roleCache) SetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRoles, data *response.ApiResponsePaginationRoleDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(roleTrashedCacheKey, req.Page, req.PageSize, req.Search)

	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *roleCache) GetCachedRoles(ctx context.Context, req *requests.FindAllRoles) (*response.ApiResponsePaginationRole, bool) {
	key := fmt.Sprintf(roleAllCacheKey, req.Page, req.PageSize, req.Search)

	payload, found := GetFromCache[response.ApiResponsePaginationRole](ctx, s.store, key)

	if !found || payload.Data == nil {
		return nil, false
	}

	return payload, true
}

func (s *roleCache) GetCachedRoleByUserId(ctx context.Context, userId int) (*response.ApiResponsesRole, bool) {
	key := fmt.Sprintf(roleByUserIdCacheKey, userId)

	result, found := GetFromCache[*response.ApiResponsesRole](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *roleCache) GetCachedRoleById(ctx context.Context, id int) (*response.ApiResponseRole, bool) {
	key := fmt.Sprintf(roleByIdCacheKey, id)

	result, found := GetFromCache[response.ApiResponseRole](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *roleCache) GetCachedRoleActive(ctx context.Context, req *requests.FindAllRoles) (*response.ApiResponsePaginationRoleDeleteAt, bool) {
	key := fmt.Sprintf(roleActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload, found := GetFromCache[response.ApiResponsePaginationRoleDeleteAt](ctx, s.store, key)

	if !found || payload.Data == nil {
		return nil, false
	}

	return payload, true
}

func (s *roleCache) GetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRoles) (*response.ApiResponsePaginationRoleDeleteAt, bool) {
	key := fmt.Sprintf(roleTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload, found := GetFromCache[response.ApiResponsePaginationRoleDeleteAt](ctx, s.store, key)

	if !found || payload.Data == nil {
		return nil, false
	}

	return payload, true
}
