package mapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/domain/response"
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb/auth"
)

type authResponseMapper struct {
}

func NewAuthResponseMapper() AuthResponseMapper {
	return &authResponseMapper{}
}

func (s *authResponseMapper) ToResponseVerifyCode(res *pb.ApiResponseVerifyCode) *response.ApiResponseVerifyCode {
	return &response.ApiResponseVerifyCode{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (s *authResponseMapper) ToResponseForgotPassword(res *pb.ApiResponseForgotPassword) *response.ApiResponseForgotPassword {
	return &response.ApiResponseForgotPassword{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (s *authResponseMapper) ToResponseResetPassword(res *pb.ApiResponseResetPassword) *response.ApiResponseResetPassword {
	return &response.ApiResponseResetPassword{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (s *authResponseMapper) ToResponseLogin(res *pb.ApiResponseLogin) *response.ApiResponseLogin {
	if res == nil {
		return &response.ApiResponseLogin{
			Status:  "error",
			Message: "response is nil",
			Data:    nil,
		}
	}

	var tokenResponse *response.TokenResponse
	if res.Data != nil {
		tokenResponse = &response.TokenResponse{
			AccessToken:  res.Data.AccessToken,
			RefreshToken: res.Data.RefreshToken,
		}
	}

	return &response.ApiResponseLogin{
		Status:  res.Status,
		Message: res.Message,
		Data:    tokenResponse,
	}
}

func (s *authResponseMapper) ToResponseRegister(res *pb.ApiResponseRegister) *response.ApiResponseRegister {
	if res == nil {
		return &response.ApiResponseRegister{
			Status:  "error",
			Message: "response is nil",
			Data:    nil,
		}
	}

	var userResponse *response.UserResponse
	if res.Data != nil {
		userResponse = &response.UserResponse{
			ID:        int(res.Data.Id),
			FirstName: res.Data.Firstname,
			LastName:  res.Data.Lastname,
			Email:     res.Data.Email,
			CreatedAt: res.Data.CreatedAt,
			UpdatedAt: res.Data.UpdatedAt,
		}
	}

	return &response.ApiResponseRegister{
		Status:  res.Status,
		Message: res.Message,
		Data:    userResponse,
	}
}

func (s *authResponseMapper) ToResponseRefreshToken(res *pb.ApiResponseRefreshToken) *response.ApiResponseRefreshToken {
	if res == nil {
		return &response.ApiResponseRefreshToken{
			Status:  "error",
			Message: "response is nil",
			Data:    nil,
		}
	}

	var tokenResponse *response.TokenResponse
	if res.Data != nil {
		tokenResponse = &response.TokenResponse{
			AccessToken:  res.Data.AccessToken,
			RefreshToken: res.Data.RefreshToken,
		}
	}

	return &response.ApiResponseRefreshToken{
		Status:  res.Status,
		Message: res.Message,
		Data:    tokenResponse,
	}
}

func (s *authResponseMapper) ToResponseGetMe(res *pb.ApiResponseGetMe) *response.ApiResponseGetMe {
	if res == nil {
		return &response.ApiResponseGetMe{
			Status:  "error",
			Message: "response is nil",
			Data:    nil,
		}
	}

	var userResponse *response.UserResponse
	if res.Data != nil {
		userResponse = &response.UserResponse{
			ID:        int(res.Data.Id),
			FirstName: res.Data.Firstname,
			LastName:  res.Data.Lastname,
			Email:     res.Data.Email,
			CreatedAt: res.Data.CreatedAt,
			UpdatedAt: res.Data.UpdatedAt,
		}
	}

	return &response.ApiResponseGetMe{
		Status:  res.Status,
		Message: res.Message,
		Data:    userResponse,
	}
}
