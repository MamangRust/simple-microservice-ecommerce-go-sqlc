package grpclientmapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/response"
	pbrole "github.com/MamangRust/simple_microservice_ecommerce_pb/role"
)

type RoleClientResponseMapper interface {
	ToApiResponseRole(pbResponse *pbrole.ApiResponseRole) *response.ApiResponseRole
}

type roleClientResponseMapper struct {
}

func NewRoleClientResponseMapper() RoleClientResponseMapper {
	return &roleClientResponseMapper{}
}

func (m *roleClientResponseMapper) ToResponseRole(role *pbrole.RoleResponse) *response.RoleResponse {
	return &response.RoleResponse{
		ID:        int(role.Id),
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}
}

func (m *roleClientResponseMapper) ToApiResponseRole(pbResponse *pbrole.ApiResponseRole) *response.ApiResponseRole {
	return &response.ApiResponseRole{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseRole(pbResponse.Data),
	}
}
