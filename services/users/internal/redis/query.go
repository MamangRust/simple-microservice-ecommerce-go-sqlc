package mencache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/response"
)

const (
	userAllCacheKey     = "user:all:page:%d:pageSize:%d:search:%s"
	userActiveCacheKey  = "user:active:page:%d:pageSize:%d:search:%s"
	userTrashedCacheKey = "user:trashed:page:%d:pageSize:%d:search:%s"
	userByIdCacheKey    = "user:id:%d"
	userByEmailCacheKey = "user:email:%s"
	ttlDefault          = 5 * time.Minute
)

type userCachedResponse struct {
	Data         []*response.UserResponse `json:"data"`
	TotalRecords *int                     `json:"total_records"`
}

type userCachedResponseDeleteAt struct {
	Data         []*response.UserResponseDeleteAt `json:"data"`
	TotalRecords *int                             `json:"total_records"`
}

type userQueryCache struct {
	store *CacheStore
}

func NewUserQueryCache(store *CacheStore) UserQueryCache {
	return &userQueryCache{store: store}
}

func (m *userQueryCache) SetCachedUsers(ctx context.Context, req *requests.FindAllUsers, data []*response.UserResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.UserResponse{}
	}

	key := fmt.Sprintf(userAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &userCachedResponse{Data: data, TotalRecords: total}
	SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *userQueryCache) SetCachedUserById(ctx context.Context, data *response.UserResponse) {
	if data == nil {
		data = &response.UserResponse{}
	}

	key := fmt.Sprintf(userByIdCacheKey, data.ID)
	SetToCache(ctx, m.store, key, data, ttlDefault)
}

func (m *userQueryCache) SetCachedUserByEmail(ctx context.Context, data *response.UserResponse) {
	if data == nil {
		data = &response.UserResponse{}
	}

	key := fmt.Sprintf(userByEmailCacheKey, data.Email)
	SetToCache(ctx, m.store, key, data, ttlDefault)
}

func (m *userQueryCache) SetCachedUserActive(ctx context.Context, req *requests.FindAllUsers, data []*response.UserResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.UserResponseDeleteAt{}
	}

	key := fmt.Sprintf(userActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload := &userCachedResponseDeleteAt{Data: data, TotalRecords: total}
	SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *userQueryCache) SetCachedUserTrashed(ctx context.Context, req *requests.FindAllUsers, data []*response.UserResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.UserResponseDeleteAt{}
	}

	key := fmt.Sprintf(userTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &userCachedResponseDeleteAt{Data: data, TotalRecords: total}
	SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *userQueryCache) GetCachedUsers(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponse, *int, bool) {
	key := fmt.Sprintf(userAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[userCachedResponse](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *userQueryCache) GetCachedUserById(ctx context.Context, id int) (*response.UserResponse, bool) {
	key := fmt.Sprintf(userByIdCacheKey, id)

	result, found := GetFromCache[*response.UserResponse](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *userQueryCache) GetCachedUserByEmail(ctx context.Context, email string) (*response.UserResponse, bool) {
	key := fmt.Sprintf(userByEmailCacheKey, email)

	result, found := GetFromCache[*response.UserResponse](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *userQueryCache) GetCachedUserByUserId(ctx context.Context, userId int) ([]*response.UserResponse, bool) {
	key := fmt.Sprintf(userByIdCacheKey, userId)

	result, found := GetFromCache[[]*response.UserResponse](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *userQueryCache) GetCachedUserActive(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(userActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[userCachedResponseDeleteAt](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *userQueryCache) GetCachedUserTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(userTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[userCachedResponseDeleteAt](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}
