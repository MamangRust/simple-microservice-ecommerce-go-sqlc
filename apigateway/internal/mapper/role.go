package mapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/domain/response"
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb/role"
)

type roleResponseMapper struct {
}

func NewRoleResponseMapper() RoleResponseMapper {
	return &roleResponseMapper{}
}

func (s *roleResponseMapper) ToApiResponseRoleAll(pbResponse *pb.ApiResponseRoleAll) *response.ApiResponseRoleAll {
	return &response.ApiResponseRoleAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (s *roleResponseMapper) ToApiResponseRoleDelete(pbResponse *pb.ApiResponseRoleDelete) *response.ApiResponseRoleDelete {
	return &response.ApiResponseRoleDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (s *roleResponseMapper) ToApiResponseRole(pbResponse *pb.ApiResponseRole) *response.ApiResponseRole {
	return &response.ApiResponseRole{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    s.mapResponseRole(pbResponse.Data),
	}
}

func (s *roleResponseMapper) ToApiResponseRoleDeleteAt(pbResponse *pb.ApiResponseRoleDeleteAt) *response.ApiResponseRoleDeleteAt {
	return &response.ApiResponseRoleDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    s.mapResponseRoleDeleteAt(pbResponse.Data),
	}
}

func (s *roleResponseMapper) ToApiResponsesRole(pbResponse *pb.ApiResponsesRole) *response.ApiResponsesRole {
	return &response.ApiResponsesRole{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    s.mapResponsesRole(pbResponse.Data),
	}
}

func (s *roleResponseMapper) ToApiResponsePaginationRole(pbResponse *pb.ApiResponsePaginationRole) *response.ApiResponsePaginationRole {
	return &response.ApiResponsePaginationRole{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       s.mapResponsesRole(pbResponse.Data),
		Pagination: mapPaginationMeta(pbResponse.PaginationMeta),
	}
}

func (s *roleResponseMapper) ToApiResponsePaginationRoleDeleteAt(pbResponse *pb.ApiResponsePaginationRoleDeleteAt) *response.ApiResponsePaginationRoleDeleteAt {
	return &response.ApiResponsePaginationRoleDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       s.mapResponsesRoleDeleteAt(pbResponse.Data),
		Pagination: mapPaginationMeta(pbResponse.PaginationMeta),
	}
}

func (s *roleResponseMapper) mapResponseRole(role *pb.RoleResponse) *response.RoleResponse {
	return &response.RoleResponse{
		ID:        int(role.Id),
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}
}

func (s *roleResponseMapper) mapResponsesRole(roles []*pb.RoleResponse) []*response.RoleResponse {
	var responseRoles []*response.RoleResponse

	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapResponseRole(role))
	}

	return responseRoles
}

func (s *roleResponseMapper) mapResponseRoleDeleteAt(role *pb.RoleResponseDeleteAt) *response.RoleResponseDeleteAt {
	var deletedAt string
	if role.DeletedAt != nil {
		deletedAt = role.DeletedAt.Value
	}

	return &response.RoleResponseDeleteAt{
		ID:        int(role.Id),
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
		DeletedAt: &deletedAt,
	}
}

func (s *roleResponseMapper) mapResponsesRoleDeleteAt(roles []*pb.RoleResponseDeleteAt) []*response.RoleResponseDeleteAt {
	var responseRoles []*response.RoleResponseDeleteAt

	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapResponseRoleDeleteAt(role))
	}

	return responseRoles
}
