package userroleprotomapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/response"
	pbuserrole "github.com/MamangRust/simple_microservice_ecommerce_pb/user_role"
)

type UserRoleProtoMappper interface {
	ToProtoResponseUserRole(status string, message string, userRole *response.UserRoleResponse) *pbuserrole.ApiResponseUserRole
}
