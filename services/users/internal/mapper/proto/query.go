package protomapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/response"
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb"
	pbuser "github.com/MamangRust/simple_microservice_ecommerce_pb/user"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type userQueryProtoMapper struct {
}

func NewUserQueryProtoMapper() UserQueryProtoMapper {
	return &userQueryProtoMapper{}
}

func (u *userQueryProtoMapper) ToProtoResponseUser(status string, message string, pbResponse *response.UserResponse) *pbuser.ApiResponseUser {
	return &pbuser.ApiResponseUser{
		Status:  status,
		Message: message,
		Data:    u.mapResponseUser(pbResponse),
	}
}

func (u *userQueryProtoMapper) ToProtoResponseUserWithPassword(status string, message string, pbResponse *response.UserResponseWithPassword) *pbuser.ApiResponseUserWithPassword {
	return &pbuser.ApiResponseUserWithPassword{
		Status:  status,
		Message: message,
		Data:    u.mapResponseUserPassword(pbResponse),
	}
}

func (u *userQueryProtoMapper) ToProtoResponseUserDeleteAt(status string, message string, pbResponse *response.UserResponseDeleteAt) *pbuser.ApiResponseUserDeleteAt {
	return &pbuser.ApiResponseUserDeleteAt{
		Status:  status,
		Message: message,
		Data:    u.mapResponseUserDeleteAt(pbResponse),
	}
}

func (u *userQueryProtoMapper) ToProtoResponsePaginationUser(pagination *pb.Pagination, status string, message string, users []*response.UserResponse) *pbuser.ApiResponsePaginationUser {
	return &pbuser.ApiResponsePaginationUser{
		Status:     status,
		Message:    message,
		Data:       u.mapResponsesUser(users),
		Pagination: MapPaginationMeta(pagination),
	}
}

func (u *userQueryProtoMapper) ToProtoResponsePaginationUserDeleteAt(pagination *pb.Pagination, status string, message string, users []*response.UserResponseDeleteAt) *pbuser.ApiResponsePaginationUserDeleteAt {
	return &pbuser.ApiResponsePaginationUserDeleteAt{
		Status:     status,
		Message:    message,
		Data:       u.mapResponsesUserDeleteAt(users),
		Pagination: MapPaginationMeta(pagination),
	}
}

func (u *userQueryProtoMapper) mapResponseUser(user *response.UserResponse) *pbuser.UserResponse {
	return &pbuser.UserResponse{
		Id:        int32(user.ID),
		Firstname: user.FirstName,
		Lastname:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (u *userQueryProtoMapper) mapResponsesUser(users []*response.UserResponse) []*pbuser.UserResponse {
	var mappedUsers []*pbuser.UserResponse

	for _, user := range users {
		mappedUsers = append(mappedUsers, u.mapResponseUser(user))
	}

	return mappedUsers
}

func (u *userQueryProtoMapper) mapResponseUserPassword(user *response.UserResponseWithPassword) *pbuser.UserResponseWithPassword {
	return &pbuser.UserResponseWithPassword{
		Id:        int32(user.ID),
		Firstname: user.FirstName,
		Lastname:  user.LastName,
		Password:  user.Password,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (u *userQueryProtoMapper) mapResponseUserDeleteAt(user *response.UserResponseDeleteAt) *pbuser.UserResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if user.DeletedAt != nil {
		deletedAt = wrapperspb.String(*user.DeletedAt)
	}

	return &pbuser.UserResponseDeleteAt{
		Id:        int32(user.ID),
		Firstname: user.FirstName,
		Lastname:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func (u *userQueryProtoMapper) mapResponsesUserDeleteAt(users []*response.UserResponseDeleteAt) []*pbuser.UserResponseDeleteAt {
	var mappedUsers []*pbuser.UserResponseDeleteAt

	for _, user := range users {
		mappedUsers = append(mappedUsers, u.mapResponseUserDeleteAt(user))
	}

	return mappedUsers
}
