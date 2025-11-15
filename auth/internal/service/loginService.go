package service

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/response"
	grpcclient "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/grpc_client"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/hash"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/logger"
	"go.uber.org/zap"
)

type LoginServiceDeps struct {
	hash         hash.HashPassword
	userClient   grpcclient.UserGrpcClientHandler
	tokenService *tokenService
	logger       logger.LoggerInterface
}

type loginService struct {
	hash         hash.HashPassword
	userClient   grpcclient.UserGrpcClientHandler
	tokenService *tokenService
	logger       logger.LoggerInterface
}

func NewLoginService(params *LoginServiceDeps) LoginService {
	return &loginService{
		hash:         params.hash,
		userClient:   params.userClient,
		logger:       params.logger,
		tokenService: params.tokenService,
	}
}

func (s *loginService) Login(ctx context.Context, request *requests.AuthRequest) (*response.TokenResponse, *response.ErrorResponse) {
	s.logger.Info("Login attempt", zap.String("email", request.Email))

	res, err := s.userClient.FindByEmailAndVerify(ctx, request.Email)

	if err != nil {
		s.logger.Error("User not found or not verified", zap.String("email", request.Email), zap.String("error_message", err.Message))
		return nil, response.NewErrorResponse("user not found or not verified", 404)
	}

	resPb := res.Data

	s.logger.Info("info password", zap.String("password", resPb.Password))

	errResp := s.hash.ComparePassword(resPb.Password, request.Password)

	if errResp != nil {
		s.logger.Error("Invalid password", zap.String("email", request.Email), zap.Error(errResp))
		return nil, response.NewErrorResponse("invalid email or password", 401)
	}

	accessToken, errResp := s.tokenService.createAccessToken(resPb.ID)

	if errResp != nil {
		s.logger.Error("Failed to generate access token", zap.Int("user_id", resPb.ID), zap.Error(errResp))
		return nil, response.NewErrorResponse("failed generate access token", 400)
	}

	refreshToken, errResp := s.tokenService.createRefreshToken(ctx, resPb.ID)

	if errResp != nil {
		s.logger.Error("Failed to generate refresh token", zap.Int("user_id", resPb.ID), zap.Error(errResp))
		return nil, response.NewErrorResponse("failed generate refresh token", 400)
	}

	tokenResp := &response.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	s.logger.Info("Login successful", zap.Int("user_id", resPb.ID), zap.String("email", request.Email))

	return tokenResp, nil
}
