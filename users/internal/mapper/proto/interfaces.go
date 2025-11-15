package protomapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/response"
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb"
	pbuser "github.com/MamangRust/simple_microservice_ecommerce_pb/user"
)

type UserBaseProtoMapper interface {
	ToProtoResponseUser(status string, message string, pbResponse *response.UserResponse) *pbuser.ApiResponseUser
}

type UserQueryProtoMapper interface {
	ToProtoResponseUser(status string, message string, pbResponse *response.UserResponse) *pbuser.ApiResponseUser
	ToProtoResponseUserWithPassword(status string, message string, pbResponse *response.UserResponseWithPassword) *pbuser.ApiResponseUserWithPassword
	ToProtoResponseUserDeleteAt(status string, message string, pbResponse *response.UserResponseDeleteAt) *pbuser.ApiResponseUserDeleteAt

	ToProtoResponsePaginationUserDeleteAt(pagination *pb.Pagination, status string, message string, users []*response.UserResponseDeleteAt) *pbuser.ApiResponsePaginationUserDeleteAt
	ToProtoResponsePaginationUser(pagination *pb.Pagination, status string, message string, users []*response.UserResponse) *pbuser.ApiResponsePaginationUser
}

type UserCommandProtoMapper interface {
	ToProtoResponseUser(status string, message string, pbResponse *response.UserResponse) *pbuser.ApiResponseUser
	ToProtoResponseUserWithPassword(status string, message string, pbResponse *response.UserResponseWithPassword) *pbuser.ApiResponseUserWithPassword
	ToProtoResponseUserDeleteAt(status string, message string, pbResponse *response.UserResponseDeleteAt) *pbuser.ApiResponseUserDeleteAt
	ToProtoResponseUserDelete(status string, message string) *pbuser.ApiResponseUserDelete
	ToProtoResponseUserAll(status string, message string) *pbuser.ApiResponseUserAll
}
