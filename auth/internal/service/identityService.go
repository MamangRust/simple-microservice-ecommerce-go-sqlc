package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/response"
	grpcclient "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/grpc_client"
	refreshtokenrepository "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/repository/refresh_token"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/auth"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/logger"
	"go.uber.org/zap"
)

type IdentityServiceDeps struct {
	auth         auth.TokenManager
	tokenService *tokenService
	user         grpcclient.UserGrpcClientHandler
	refreshToken refreshtokenrepository.RefreshTokenCommandRepository
	logger       logger.LoggerInterface
}

type identityService struct {
	auth         auth.TokenManager
	tokenService *tokenService
	user         grpcclient.UserGrpcClientHandler
	refreshtoken refreshtokenrepository.RefreshTokenCommandRepository
	logger       logger.LoggerInterface
}

func NewIdentityService(params *IdentityServiceDeps) *identityService {
	return &identityService{
		auth:         params.auth,
		tokenService: params.tokenService,
		user:         params.user,
		refreshtoken: params.refreshToken,
		logger:       params.logger,
	}
}

func (s *identityService) RefreshToken(ctx context.Context, token string) (*response.TokenResponse, *response.ErrorResponse) {
	s.logger.Info("Token refresh attempt")

	userIdStr, err := s.auth.ValidateToken(token)

	if err != nil {
		s.logger.Error("Token validation failed", zap.Error(err))
		if errors.Is(err, auth.ErrTokenExpired) {
			return nil, response.NewErrorResponse("refresh token expired", 401)
		}

		return nil, response.NewErrorResponse(err.Error(), 401)
	}

	if err := s.refreshtoken.DeleteRefreshToken(ctx, token); err != nil {
		s.logger.Error("Failed to delete refresh token", zap.String("token", token), zap.Error(err))
		return nil, &response.ErrorResponse{}
	}

	userId, err := strconv.Atoi(userIdStr)

	if err != nil {
		s.logger.Error("Invalid user ID format", zap.String("user_id_str", userIdStr), zap.Error(err))
		return nil, response.NewErrorResponse("invalid user id format", 400)
	}

	accessToken, err := s.tokenService.createAccessToken(userId)

	if err != nil {
		s.logger.Error("Failed to create access token", zap.Int("user_id", userId), zap.Error(err))
		return nil, response.NewErrorResponse("failed to create access token", 400)
	}

	refreshToken, err := s.tokenService.createRefreshToken(ctx, userId)

	if err != nil {
		s.logger.Error("Failed to create refresh token", zap.Int("user_id", userId), zap.Error(err))
		return nil, response.NewErrorResponse("failed to update refresh token", 400)
	}

	expiryTime := time.Now().Add(24 * time.Hour)
	updateRequest := &requests.UpdateRefreshToken{
		UserId:    userId,
		Token:     refreshToken,
		ExpiresAt: expiryTime.Format("2006-01-02 15:04:05"),
	}

	if _, err = s.refreshtoken.UpdateRefreshToken(ctx, updateRequest); err != nil {
		s.logger.Error("Failed to update refresh token", zap.Int("user_id", userId), zap.Error(err))
		return nil, response.NewErrorResponse("failed to update refresh token", 500)
	}

	s.logger.Info("Token refresh successful", zap.Int("user_id", userId))

	return &response.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *identityService) GetMe(ctx context.Context, userid int) (*response.UserResponse, *response.ErrorResponse) {
	s.logger.Info("Get user profile request", zap.Int("user_id", userid))

	user, err := s.user.FindById(ctx, int32(userid))

	if err != nil {
		s.logger.Error("User not found", zap.Int("user_id", userid), zap.Any("error", err))
		return nil, response.NewErrorResponse("user not found", 404)
	}

	s.logger.Info("Get user profile successful", zap.Int("user_id", userid))

	return user.Data, nil
}
