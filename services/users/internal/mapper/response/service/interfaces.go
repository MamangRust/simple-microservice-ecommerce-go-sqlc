package responsemapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/response"
)

type UserBaseResponseMapper interface {
	ToUserResponse(user *record.UserRecord) *response.UserResponse
}

type UserQueryResponseMapper interface {
	UserBaseResponseMapper

	ToUserWithPasswordResponse(user *record.UserRecord) *response.UserResponseWithPassword

	ToUsersResponse(users []*record.UserRecord) []*response.UserResponse

	ToUsersResponseDeleteAt(users []*record.UserRecord) []*response.UserResponseDeleteAt
}

type UserCommandResponseMapper interface {
	UserBaseResponseMapper

	ToUserResponseDeleteAt(user *record.UserRecord) *response.UserResponseDeleteAt
}
