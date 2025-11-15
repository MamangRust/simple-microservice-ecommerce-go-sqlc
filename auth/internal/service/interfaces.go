package service

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/response"
)

type RegistrationService interface {
	Register(ctx context.Context, request *requests.RegisterRequest) (*response.UserResponse, *response.ErrorResponse)
}

type LoginService interface {
	Login(ctx context.Context, request *requests.AuthRequest) (*response.TokenResponse, *response.ErrorResponse)
}

type PasswordResetService interface {
	ForgotPassword(ctx context.Context, email string) (bool, *response.ErrorResponse)
	ResetPassword(ctx context.Context, request *requests.CreateResetPasswordRequest) (bool, *response.ErrorResponse)
	VerifyCode(ctx context.Context, code string) (bool, *response.ErrorResponse)
}

type IdentifyService interface {
	RefreshToken(ctx context.Context, token string) (*response.TokenResponse, *response.ErrorResponse)
	GetMe(ctx context.Context, userid int) (*response.UserResponse, *response.ErrorResponse)
}
