package roleprotomapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/response"
	helperproto "github.com/MamangRust/simple_microservice_ecommerce/role/internal/mapper/proto/helpers"
	pbrole "github.com/MamangRust/simple_microservice_ecommerce_pb/role"
)

type roleCommandProtoMapper struct {
}

func NewRoleCommandProtoMapper() RoleCommandProtoMapper {
	return &roleCommandProtoMapper{}
}

func (s *roleCommandProtoMapper) ToProtoResponseRole(status string, message string, pbroleResponse *response.RoleResponse) *pbrole.ApiResponseRole {
	return &pbrole.ApiResponseRole{
		Status:  status,
		Message: message,
		Data:    s.mapResponseRole(pbroleResponse),
	}
}

func (s *roleCommandProtoMapper) ToProtoResponseRoleDeleteAt(status string, message string, pbroleResponse *response.RoleResponseDeleteAt) *pbrole.ApiResponseRoleDeleteAt {
	return &pbrole.ApiResponseRoleDeleteAt{
		Status:  status,
		Message: message,
		Data:    s.mapResponseRoleDeleteAt(pbroleResponse),
	}
}

func (s *roleCommandProtoMapper) ToProtoResponseRoleAll(status string, message string) *pbrole.ApiResponseRoleAll {
	return &pbrole.ApiResponseRoleAll{
		Status:  status,
		Message: message,
	}
}

func (s *roleCommandProtoMapper) ToProtoResponseRoleDelete(status string, message string) *pbrole.ApiResponseRoleDelete {
	return &pbrole.ApiResponseRoleDelete{
		Status:  status,
		Message: message,
	}
}

func (s *roleCommandProtoMapper) mapResponseRole(role *response.RoleResponse) *pbrole.RoleResponse {
	return &pbrole.RoleResponse{
		Id:        int32(role.ID),
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}
}

func (s *roleCommandProtoMapper) mapResponseRoleDeleteAt(role *response.RoleResponseDeleteAt) *pbrole.RoleResponseDeleteAt {
	res := &pbrole.RoleResponseDeleteAt{
		Id:        int32(role.ID),
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}

	if role.DeletedAt != nil {
		res.DeletedAt = helperproto.StringPtrToWrapper(role.DeletedAt)
	}

	return res
}
