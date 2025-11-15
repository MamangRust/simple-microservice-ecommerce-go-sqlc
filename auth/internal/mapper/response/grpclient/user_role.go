package grpclientmapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/response"
	pbuserrole "github.com/MamangRust/simple_microservice_ecommerce_pb/user_role"
)

type UserRoleClientResponseMapper interface {
	ToApiResponseUserRole(pbResponse *pbuserrole.ApiResponseUserRole) *response.ApiResponseUserRole
}

type userRoleClientResponseMapper struct {
}

func NewUserRoleClientResponseMapper() UserRoleClientResponseMapper {
	return &userRoleClientResponseMapper{}
}

func (u *userRoleClientResponseMapper) ToResponseUserRole(userRole *pbuserrole.UserRoleResponse) *response.UserRoleResponse {
	return &response.UserRoleResponse{
		UserId: int(userRole.Userid),
		RoleId: int(userRole.Roleid),
	}
}

func (u *userRoleClientResponseMapper) ToApiResponseUserRole(pbResponse *pbuserrole.ApiResponseUserRole) *response.ApiResponseUserRole {
	return &response.ApiResponseUserRole{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    u.ToResponseUserRole(pbResponse.Data),
	}
}
