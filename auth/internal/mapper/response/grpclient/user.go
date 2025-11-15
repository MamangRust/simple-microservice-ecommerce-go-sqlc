package grpclientmapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/response"
	pbuser "github.com/MamangRust/simple_microservice_ecommerce_pb/user"
)

type UserClientResponseMapper interface {
	ToApiResponseUser(pbResponse *pbuser.ApiResponseUser) *response.ApiResponseUser
	ToApiResponseUserWithPassword(pbResponse *pbuser.ApiResponseUserWithPassword) *response.ApiResponseUserWithPassword
}

type userClientResponseMapper struct {
}

func NewUserClientResponseMapper() UserClientResponseMapper {
	return &userClientResponseMapper{}
}

func (u *userClientResponseMapper) ToResponseUser(user *pbuser.UserResponse) *response.UserResponse {
	return &response.UserResponse{
		ID:        int(user.Id),
		FirstName: user.Firstname,
		LastName:  user.Lastname,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (u *userClientResponseMapper) ToApiResponseUser(pbResponse *pbuser.ApiResponseUser) *response.ApiResponseUser {
	return &response.ApiResponseUser{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    u.ToResponseUser(pbResponse.Data),
	}
}

func (u *userClientResponseMapper) ToResponseUserWithPassword(user *pbuser.UserResponseWithPassword) *response.UserWithPasswordResponse {
	return &response.UserWithPasswordResponse{
		ID:        int(user.Id),
		FirstName: user.Firstname,
		LastName:  user.Lastname,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (u *userClientResponseMapper) ToApiResponseUserWithPassword(pbResponse *pbuser.ApiResponseUserWithPassword) *response.ApiResponseUserWithPassword {
	return &response.ApiResponseUserWithPassword{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    u.ToResponseUserWithPassword(pbResponse.Data),
	}
}
