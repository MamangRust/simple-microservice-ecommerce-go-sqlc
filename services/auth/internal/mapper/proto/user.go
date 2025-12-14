package protomapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/response"
	pbuser "github.com/MamangRust/simple_microservice_ecommerce_pb/user"
)

type userProtoMapper struct {
}

func NewUserProtoMapper() *userProtoMapper {
	return &userProtoMapper{}
}

func (u *userProtoMapper) ToProtoResponseUser(status string, message string, pbResponse *response.UserResponse) *pbuser.ApiResponseUser {
	return &pbuser.ApiResponseUser{
		Status:  status,
		Message: message,
		Data:    u.mapResponseUser(pbResponse),
	}
}

func (u *userProtoMapper) mapResponseUser(user *response.UserResponse) *pbuser.UserResponse {
	return &pbuser.UserResponse{
		Id:        int32(user.ID),
		Firstname: user.FirstName,
		Lastname:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
