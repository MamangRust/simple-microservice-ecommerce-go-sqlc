package responsemapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/response"
)

type userCommandResponseMapper struct{}

func NewUserCommandResponseMapper() UserCommandResponseMapper {
	return &userCommandResponseMapper{}
}

func (s *userCommandResponseMapper) ToUserResponse(user *record.UserRecord) *response.UserResponse {
	return &response.UserResponse{
		ID:         user.ID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}

func (s *userCommandResponseMapper) ToUserResponseDeleteAt(user *record.UserRecord) *response.UserResponseDeleteAt {
	return &response.UserResponseDeleteAt{
		ID:         user.ID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
		DeletedAt:  user.DeletedAt,
	}
}
