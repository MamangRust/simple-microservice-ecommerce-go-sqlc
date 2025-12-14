package roleprotomapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/response"
	helperproto "github.com/MamangRust/simple_microservice_ecommerce/role/internal/mapper/proto/helpers"
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb"
	pbrole "github.com/MamangRust/simple_microservice_ecommerce_pb/role"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type roleQueryProtoMapper struct {
}

func NewRoleQueryProtoMapper() RoleQueryProtoMapper {
	return &roleQueryProtoMapper{}
}

func (s *roleQueryProtoMapper) ToProtoResponseRole(status string, message string, pbroleResponse *response.RoleResponse) *pbrole.ApiResponseRole {
	return &pbrole.ApiResponseRole{
		Status:  status,
		Message: message,
		Data:    s.mapResponseRole(pbroleResponse),
	}
}

func (s *roleQueryProtoMapper) ToProtoResponsesRole(status string, message string, pbroleResponse []*response.RoleResponse) *pbrole.ApiResponsesRole {
	return &pbrole.ApiResponsesRole{
		Status:  status,
		Message: message,
		Data:    s.mapResponsesRole(pbroleResponse),
	}
}

func (s *roleQueryProtoMapper) ToProtoResponsePaginationRole(pagination *pb.Pagination, status string, message string, pbroleResponse []*response.RoleResponse) *pbrole.ApiResponsePaginationRole {
	return &pbrole.ApiResponsePaginationRole{
		Status:         status,
		Message:        message,
		Data:           s.mapResponsesRole(pbroleResponse),
		PaginationMeta: MapPaginationMeta(pagination),
	}
}

func (s *roleQueryProtoMapper) ToProtoResponsePaginationRoleDeleteAt(pagination *pb.Pagination, status string, message string, pbroleResponse []*response.RoleResponseDeleteAt) *pbrole.ApiResponsePaginationRoleDeleteAt {
	return &pbrole.ApiResponsePaginationRoleDeleteAt{
		Status:         status,
		Message:        message,
		Data:           s.mapResponsesRoleDeleteAt(pbroleResponse),
		PaginationMeta: MapPaginationMeta(pagination),
	}
}

func (s *roleQueryProtoMapper) mapResponseRole(role *response.RoleResponse) *pbrole.RoleResponse {
	return &pbrole.RoleResponse{
		Id:        int32(role.ID),
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}
}

func (s *roleQueryProtoMapper) mapResponsesRole(roles []*response.RoleResponse) []*pbrole.RoleResponse {
	var responseRoles []*pbrole.RoleResponse

	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapResponseRole(role))
	}

	return responseRoles
}

func (s *roleQueryProtoMapper) mapResponseRoleDeleteAt(role *response.RoleResponseDeleteAt) *pbrole.RoleResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if role.DeletedAt != nil {
		deletedAt = helperproto.StringPtrToWrapper(role.DeletedAt)
	}

	return &pbrole.RoleResponseDeleteAt{
		Id:        int32(role.ID),
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func (s *roleQueryProtoMapper) mapResponsesRoleDeleteAt(roles []*response.RoleResponseDeleteAt) []*pbrole.RoleResponseDeleteAt {
	var responseRoles []*pbrole.RoleResponseDeleteAt

	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapResponseRoleDeleteAt(role))
	}

	return responseRoles
}

func MapPaginationMeta(s *pb.Pagination) *pb.Pagination {
	return &pb.Pagination{
		CurrentPage:  int32(s.CurrentPage),
		PageSize:     int32(s.PageSize),
		TotalPages:   int32(s.TotalPages),
		TotalRecords: int32(s.TotalRecords),
	}
}
