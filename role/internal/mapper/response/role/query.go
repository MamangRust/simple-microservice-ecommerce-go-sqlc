package roleresponsemapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/response"
)

type roleQueryResponseMapper struct{}

func NewRoleQueryResponseMapper() RoleQueryResponseMapper {
	return &roleQueryResponseMapper{}
}

func (s *roleQueryResponseMapper) ToRoleResponse(role *record.RoleRecord) *response.RoleResponse {
	return &response.RoleResponse{
		ID:        role.ID,
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}
}

func (s *roleQueryResponseMapper) ToRolesResponse(roles []*record.RoleRecord) []*response.RoleResponse {
	var responseRoles []*response.RoleResponse

	for _, role := range roles {
		responseRoles = append(responseRoles, s.ToRoleResponse(role))
	}

	return responseRoles
}

func (s *roleQueryResponseMapper) ToRoleResponseDeleteAt(role *record.RoleRecord) *response.RoleResponseDeleteAt {
	return &response.RoleResponseDeleteAt{
		ID:        role.ID,
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
		DeletedAt: role.DeletedAt,
	}
}

func (s *roleQueryResponseMapper) ToRolesResponseDeleteAt(roles []*record.RoleRecord) []*response.RoleResponseDeleteAt {
	var responseRoles []*response.RoleResponseDeleteAt

	for _, role := range roles {
		responseRoles = append(responseRoles, s.ToRoleResponseDeleteAt(role))
	}

	return responseRoles
}
