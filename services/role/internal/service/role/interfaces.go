package roleservice

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/response"
)

type RoleQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllRoles) ([]*response.RoleResponse, *int, *response.ErrorResponse)
	FindByActiveRole(ctx context.Context, req *requests.FindAllRoles) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashedRole(ctx context.Context, req *requests.FindAllRoles) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(ctx context.Context, role_id int) (*response.RoleResponse, *response.ErrorResponse)
	FindByUserId(ctx context.Context, id int) ([]*response.RoleResponse, *response.ErrorResponse)
	FindByName(ctx context.Context, name string) (*response.RoleResponse, *response.ErrorResponse)
}

type RoleCommandService interface {
	CreateRole(ctx context.Context, request *requests.CreateRoleRequest) (*response.RoleResponse, *response.ErrorResponse)
	UpdateRole(ctx context.Context, request *requests.UpdateRoleRequest) (*response.RoleResponse, *response.ErrorResponse)
	TrashedRole(ctx context.Context, role_id int) (*response.RoleResponseDeleteAt, *response.ErrorResponse)
	RestoreRole(ctx context.Context, role_id int) (*response.RoleResponseDeleteAt, *response.ErrorResponse)
	DeleteRolePermanent(ctx context.Context, role_id int) (bool, *response.ErrorResponse)
	RestoreAllRole(ctx context.Context) (bool, *response.ErrorResponse)
	DeleteAllRolePermanent(ctx context.Context) (bool, *response.ErrorResponse)
}
