package responsemapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/response"
)

type userQueryResponseMapper struct{}

func NewUserQueryResponseMapper() UserQueryResponseMapper {
	return &userQueryResponseMapper{}
}

func (s *userQueryResponseMapper) ToUserResponse(user *record.UserRecord) *response.UserResponse {
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

func (s *userQueryResponseMapper) ToUserWithPasswordResponse(user *record.UserRecord) *response.UserResponseWithPassword {
	return &response.UserResponseWithPassword{
		ID:         user.ID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
		Password:   user.Password,
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}

func (s *userQueryResponseMapper) ToUsersResponse(users []*record.UserRecord) []*response.UserResponse {
	var responses []*response.UserResponse

	for _, user := range users {
		responses = append(responses, s.ToUserResponse(user))
	}

	return responses
}

func (s *userQueryResponseMapper) ToUserResponseDeleteAt(user *record.UserRecord) *response.UserResponseDeleteAt {
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

func (s *userQueryResponseMapper) ToUsersResponseDeleteAt(users []*record.UserRecord) []*response.UserResponseDeleteAt {
	var responses []*response.UserResponseDeleteAt

	for _, user := range users {
		responses = append(responses, s.ToUserResponseDeleteAt(user))
	}

	return responses
}
