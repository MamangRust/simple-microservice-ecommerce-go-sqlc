package service

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/response"
)

type UserQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponse, *int, *response.ErrorResponse)
	FindByID(ctx context.Context, id int) (*response.UserResponse, *response.ErrorResponse)
	FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse)

	FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse)

	FindByEmail(ctx context.Context, email string) (*response.UserResponse, *response.ErrorResponse)
	FindByEmailAndVerify(ctx context.Context, email string) (*response.UserResponseWithPassword, *response.ErrorResponse)
	FindByVerificationCode(ctx context.Context, code string) (*response.UserResponse, *response.ErrorResponse)
}

type UserCommandService interface {
	CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*response.UserResponse, *response.ErrorResponse)
	UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*response.UserResponse, *response.ErrorResponse)
	UpdateUserIsVerified(ctx context.Context, request *requests.UpdateUserVerifiedRequest) (*response.UserResponse, *response.ErrorResponse)
	UpdateUserPassword(ctx context.Context, request *requests.UpdateUserPasswordRequest) (*response.UserResponse, *response.ErrorResponse)
	TrashedUser(ctx context.Context, user_id int) (*response.UserResponseDeleteAt, *response.ErrorResponse)
	RestoreUser(ctx context.Context, user_id int) (*response.UserResponseDeleteAt, *response.ErrorResponse)
	DeleteUserPermanent(ctx context.Context, user_id int) (bool, *response.ErrorResponse)
	RestoreAllUser(ctx context.Context) (bool, *response.ErrorResponse)
	DeleteAllUserPermanent(ctx context.Context) (bool, *response.ErrorResponse)
}
