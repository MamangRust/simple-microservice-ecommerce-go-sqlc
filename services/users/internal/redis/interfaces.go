package mencache

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/response"
)

type UserCommandCache interface {
	DeleteCachedUser(ctx context.Context, key string)
	InvalidateAllUsers(ctx context.Context)
	InvalidateActiveUsers(ctx context.Context)
	InvalidateTrashedUsers(ctx context.Context)
}

type UserQueryCache interface {
	SetCachedUsers(ctx context.Context, req *requests.FindAllUsers, data []*response.UserResponse, total *int)
	SetCachedUserById(ctx context.Context, data *response.UserResponse)
	SetCachedUserActive(ctx context.Context, req *requests.FindAllUsers, data []*response.UserResponseDeleteAt, total *int)
	SetCachedUserTrashed(ctx context.Context, req *requests.FindAllUsers, data []*response.UserResponseDeleteAt, total *int)
	SetCachedUserByEmail(ctx context.Context, data *response.UserResponse)

	GetCachedUsers(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponse, *int, bool)
	GetCachedUserById(ctx context.Context, id int) (*response.UserResponse, bool)
	GetCachedUserActive(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, bool)
	GetCachedUserTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, bool)
	GetCachedUserByEmail(ctx context.Context, email string) (*response.UserResponse, bool)
}
