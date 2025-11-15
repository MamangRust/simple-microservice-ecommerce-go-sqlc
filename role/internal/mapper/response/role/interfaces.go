package roleresponsemapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/response"
)

type RoleBaseResponseMapper interface {
	ToRoleResponse(role *record.RoleRecord) *response.RoleResponse
}

type RoleQueryResponseMapper interface {
	RoleBaseResponseMapper

	ToRolesResponse(roles []*record.RoleRecord) []*response.RoleResponse
	ToRolesResponseDeleteAt(roles []*record.RoleRecord) []*response.RoleResponseDeleteAt
}

type RoleCommandResponseMapper interface {
	RoleBaseResponseMapper

	ToRoleResponseDeleteAt(role *record.RoleRecord) *response.RoleResponseDeleteAt
}
