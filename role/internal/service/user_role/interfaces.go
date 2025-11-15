package userroleservice

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/response"
)

type UserRoleService interface {
	AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*response.UserRoleResponse, *response.ErrorResponse)
	UpdateRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*response.UserRoleResponse, *response.ErrorResponse)
}
