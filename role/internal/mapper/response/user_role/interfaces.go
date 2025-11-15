package userroleresponsemapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/response"
)

type UserRoleResponseMappping interface {
	ToUserRoleResponse(userRole *record.UserRoleRecord) *response.UserRoleResponse
}
