package userroleprotomapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/response"
	pbuserrole "github.com/MamangRust/simple_microservice_ecommerce_pb/user_role"
)

type userRoleProtoMapper struct {
}

func NewUserRoleProtoMapper() UserRoleProtoMappper {
	return &userRoleProtoMapper{}
}

func (u *userRoleProtoMapper) ToProtoResponseUserRole(status string, message string, userRole *response.UserRoleResponse) *pbuserrole.ApiResponseUserRole {
	return &pbuserrole.ApiResponseUserRole{
		Status:  status,
		Message: message,
		Data:    u.mapResponseUseRole(userRole),
	}
}

func (u *userRoleProtoMapper) mapResponseUseRole(userRole *response.UserRoleResponse) *pbuserrole.UserRoleResponse {
	return &pbuserrole.UserRoleResponse{
		Userid: int32(userRole.UserID),
		Roleid: int32(userRole.RoleID),
	}
}
