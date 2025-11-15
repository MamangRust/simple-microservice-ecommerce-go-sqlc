package userroleresponsemapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/response"
)

type userRoleResponseMapper struct {
}

func NewUserRoleResponseMapper() UserRoleResponseMappping {
	return &userRoleResponseMapper{}
}

func (u *userRoleResponseMapper) ToUserRoleResponse(userRole *record.UserRoleRecord) *response.UserRoleResponse {
	return &response.UserRoleResponse{
		UserID: int(userRole.UserID),
		RoleID: int(userRole.RoleID),
	}
}
