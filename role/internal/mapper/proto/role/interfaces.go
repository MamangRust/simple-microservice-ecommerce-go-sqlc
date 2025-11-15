package roleprotomapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/response"
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb"
	pbrole "github.com/MamangRust/simple_microservice_ecommerce_pb/role"
)

type RoleBaseProtoMapper interface {
	ToProtoResponseRole(status string, message string, pbResponse *response.RoleResponse) *pbrole.ApiResponseRole
}

type RoleQueryProtoMapper interface {
	RoleBaseProtoMapper

	ToProtoResponsesRole(status string, message string, pbResponse []*response.RoleResponse) *pbrole.ApiResponsesRole

	ToProtoResponsePaginationRole(pagination *pb.Pagination, status string, message string, pbResponse []*response.RoleResponse) *pbrole.ApiResponsePaginationRole

	ToProtoResponsePaginationRoleDeleteAt(pagination *pb.Pagination, status string, message string, pbResponse []*response.RoleResponseDeleteAt) *pbrole.ApiResponsePaginationRoleDeleteAt
}

type RoleCommandProtoMapper interface {
	RoleBaseProtoMapper

	ToProtoResponseRoleDeleteAt(status string, message string, pbResponse *response.RoleResponseDeleteAt) *pbrole.ApiResponseRoleDeleteAt

	ToProtoResponseRoleAll(status string, message string) *pbrole.ApiResponseRoleAll

	ToProtoResponseRoleDelete(status string, message string) *pbrole.ApiResponseRoleDelete
}
