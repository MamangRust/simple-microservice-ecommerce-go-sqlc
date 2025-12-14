package protomapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/response"

	pbuser "github.com/MamangRust/simple_microservice_ecommerce_pb/user"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type userCommandProtoMapper struct {
}

func NewUserCommandProtoMapper() UserCommandProtoMapper {
	return &userCommandProtoMapper{}
}

func (u *userCommandProtoMapper) ToProtoResponseUser(status string, message string, pbResponse *response.UserResponse) *pbuser.ApiResponseUser {
	return &pbuser.ApiResponseUser{
		Status:  status,
		Message: message,
		Data:    u.mapResponseUser(pbResponse),
	}
}

func (u *userCommandProtoMapper) ToProtoResponseUserWithPassword(status string, message string, pbResponse *response.UserResponseWithPassword) *pbuser.ApiResponseUserWithPassword {
	return &pbuser.ApiResponseUserWithPassword{
		Status:  status,
		Message: message,
		Data:    u.mapResponseUserPassword(pbResponse),
	}
}

func (u *userCommandProtoMapper) ToProtoResponseUserDeleteAt(status string, message string, pbResponse *response.UserResponseDeleteAt) *pbuser.ApiResponseUserDeleteAt {
	return &pbuser.ApiResponseUserDeleteAt{
		Status:  status,
		Message: message,
		Data:    u.mapResponseUserDeleteAt(pbResponse),
	}
}

func (u *userCommandProtoMapper) ToProtoResponseUserDelete(status string, message string) *pbuser.ApiResponseUserDelete {
	return &pbuser.ApiResponseUserDelete{
		Status:  status,
		Message: message,
	}
}

func (u *userCommandProtoMapper) ToProtoResponseUserAll(status string, message string) *pbuser.ApiResponseUserAll {
	return &pbuser.ApiResponseUserAll{
		Status:  status,
		Message: message,
	}
}

func (u *userCommandProtoMapper) mapResponseUser(user *response.UserResponse) *pbuser.UserResponse {
	return &pbuser.UserResponse{
		Id:        int32(user.ID),
		Firstname: user.FirstName,
		Lastname:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (u *userCommandProtoMapper) mapResponseUserPassword(user *response.UserResponseWithPassword) *pbuser.UserResponseWithPassword {
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

func (u *userCommandProtoMapper) mapResponseUserDeleteAt(user *response.UserResponseDeleteAt) *pbuser.UserResponseDeleteAt {
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
