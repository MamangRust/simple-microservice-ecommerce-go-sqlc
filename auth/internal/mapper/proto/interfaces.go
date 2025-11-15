package protomapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/response"
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb/auth"
	pbuser "github.com/MamangRust/simple_microservice_ecommerce_pb/user"
)

type AuthProtoMapper interface {
	ToProtoResponseVerifyCode(status string, message string) *pb.ApiResponseVerifyCode
	ToProtoResponseForgotPassword(status string, message string) *pb.ApiResponseForgotPassword
	ToProtoResponseResetPassword(status string, message string) *pb.ApiResponseResetPassword
	ToProtoResponseLogin(status string, message string, response *response.TokenResponse) *pb.ApiResponseLogin
	ToProtoResponseRegister(status string, message string, response *response.UserResponse) *pb.ApiResponseRegister
	ToProtoResponseRefreshToken(status string, message string, response *response.TokenResponse) *pb.ApiResponseRefreshToken
	ToProtoResponseGetMe(status string, message string, response *response.UserResponse) *pb.ApiResponseGetMe
}

type UserProtoMapper interface {
	ToProtoResponseUser(status string, message string, pbResponse *response.UserResponse) *pbuser.ApiResponseUser
}
