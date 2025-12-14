package roleresponsemapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/response"
)

type roleCommandResponseMapper struct {
}

func NewRoleCommandResponseMapper() RoleCommandResponseMapper {
	return &roleCommandResponseMapper{}
}

func (s *roleCommandResponseMapper) ToRoleResponse(role *record.RoleRecord) *response.RoleResponse {
	return &response.RoleResponse{
		ID:        role.ID,
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}
}

func (s *roleCommandResponseMapper) ToRoleResponseDeleteAt(role *record.RoleRecord) *response.RoleResponseDeleteAt {
	return &response.RoleResponseDeleteAt{
		ID:        role.ID,
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
		DeletedAt: role.DeletedAt,
	}
}
