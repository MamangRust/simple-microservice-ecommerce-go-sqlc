package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/response"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/errorhandler"
	grpcclient "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/grpc_client"
	mencache "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/redis"
	refreshtokenrepository "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/repository/refresh_token"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/auth"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/logger"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type IdentityServiceDeps struct {
	auth         auth.TokenManager
	tokenService *tokenService
	user         grpcclient.UserGrpcClientHandler
	refreshToken refreshtokenrepository.RefreshTokenCommandRepository
	logger       logger.LoggerInterface
	errorhandler errorhandler.IdentityErrorHandler
	errorToken   errorhandler.TokenErrorHandler
	mencache     mencache.IdentityCache
}

type identityService struct {
	auth          auth.TokenManager
	tokenService  *tokenService
	user          grpcclient.UserGrpcClientHandler
	refreshtoken  refreshtokenrepository.RefreshTokenCommandRepository
	logger        logger.LoggerInterface
	observability observability.TraceLoggerObservability
	errorhandler  errorhandler.IdentityErrorHandler
	errorToken    errorhandler.TokenErrorHandler
	mencache      mencache.IdentityCache
}

func NewIdentityService(params *IdentityServiceDeps) *identityService {
	observability, _ := observability.NewObservability("identity-service", params.logger)

	return &identityService{
		auth:          params.auth,
		tokenService:  params.tokenService,
		user:          params.user,
		refreshtoken:  params.refreshToken,
		logger:        params.logger,
		observability: observability,
		errorhandler:  params.errorhandler,
		errorToken:    params.errorToken,
		mencache:      params.mencache,
	}
}

func (s *identityService) RefreshToken(ctx context.Context, token string) (*response.TokenResponse, *response.ErrorResponse) {
	const method = "RefreshToken"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("token", token))

	defer func() {
		end(status)
	}()

	if cachedUserID, found := s.mencache.GetRefreshToken(ctx, token); found {
		userId, err := strconv.Atoi(cachedUserID)
		if err == nil {
			s.mencache.DeleteRefreshToken(ctx, token)
			s.logger.Debug("Invalidated old refresh token from cache", zap.String("token", token))

			accessToken, err := s.tokenService.createAccessToken(userId)
			if err != nil {
				status = "error"

				return s.errorToken.HandleCreateAccessTokenError(err, method, "CREATE_ACCESS_TOKEN_FAILED", span, &status, zap.Int("user.id", userId))
			}

			refreshToken, err := s.tokenService.createRefreshToken(ctx, userId)
			if err != nil {
				status = "error"

				return s.errorToken.HandleCreateRefreshTokenError(err, method, "CREATE_REFRESH_TOKEN_FAILED", span, &status, zap.Int("user.id", userId))
			}

			expiryTime := time.Now().Add(24 * time.Hour)
			expirationDuration := time.Until(expiryTime)

			s.mencache.SetRefreshToken(ctx, refreshToken, expirationDuration)
			s.logger.Debug("Stored new refresh token in cache",
				zap.String("new_token", refreshToken),
				zap.Duration("expiration", expirationDuration))

			s.logger.Debug("Refresh token refreshed successfully (cached)", zap.Int("user_id", userId))
			span.SetStatus(codes.Ok, "Token refreshed successfully from cache")

			return &response.TokenResponse{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			}, nil
		}
	}

	userIdStr, err := s.auth.ValidateToken(token)

	if err != nil {
		status = "error"
		if errors.Is(err, auth.ErrTokenExpired) {
			s.mencache.DeleteRefreshToken(ctx, token)
			if err := s.refreshtoken.DeleteRefreshToken(ctx, token); err != nil {
				return s.errorhandler.HandleDeleteRefreshTokenError(err, method, "DELETE_REFRESH_TOKEN", span, &status, zap.String("token", token))
			}

			return s.errorhandler.HandleExpiredRefreshTokenError(err, method, "TOKEN_EXPIRED", span, &status, zap.String("token", token))
		}
		return s.errorhandler.HandleInvalidTokenError(err, method, "INVALID_TOKEN", span, &status, zap.String("token", token))
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		status = "error"

		return errorhandler.HandleInvalidFormatUserIDError[*response.TokenResponse](s.logger, err, method, "INVALID_USER_ID", span, &status, zap.Int("user.id", userId))
	}

	if err := s.refreshtoken.DeleteRefreshToken(ctx, token); err != nil {
		status = "error"

		return s.errorhandler.HandleDeleteRefreshTokenError(err, method, "DELETE_REFRESH_TOKEN", span, &status, zap.String("token", token))
	}

	accessToken, err := s.tokenService.createAccessToken(userId)

	if err != nil {
		status = "error"

		return s.errorToken.HandleCreateAccessTokenError(err, method, "CREATE_ACCESS_TOKEN_FAILED", span, &status, zap.Int("user.id", userId))
	}

	refreshToken, err := s.tokenService.createRefreshToken(ctx, userId)

	if err != nil {
		status = "error"

		return s.errorToken.HandleCreateRefreshTokenError(err, method, "CREATE_REFRESH_TOKEN_FAILED", span, &status, zap.Int("user.id", userId))
	}

	expiryTime := time.Now().Add(24 * time.Hour)
	updateRequest := &requests.UpdateRefreshToken{
		UserId:    userId,
		Token:     refreshToken,
		ExpiresAt: expiryTime.Format("2006-01-02 15:04:05"),
	}

	if _, err = s.refreshtoken.UpdateRefreshToken(ctx, updateRequest); err != nil {
		s.mencache.DeleteRefreshToken(ctx, refreshToken)
		status = "error"

		return s.errorhandler.HandleUpdateRefreshTokenError(err, method, "UPDATE_REFRESH_TOKEN_FAILED", span, &status, zap.Int("user.id", userId))
	}

	expirationDuration := time.Until(expiryTime)

	s.mencache.SetRefreshToken(ctx, refreshToken, expirationDuration)
	s.mencache.DeleteRefreshToken(ctx, token)

	logSuccess("Refresh token refreshed successfully", zap.Int("user.id", userId))

	return &response.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *identityService) GetMe(ctx context.Context, userid int) (*response.UserResponse, *response.ErrorResponse) {
	const method = "GetMe"
	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("userId", userid))

	defer func() {
		end(status)
	}()

	userKey := strconv.Itoa(userid)

	if cachedUser, found := s.mencache.GetCachedUserInfo(ctx, userKey); found {
		logSuccess("User info retrieved from cache", zap.Int("user.id", userid))
		return cachedUser, nil
	}

	user, err := s.user.FindById(ctx, int32(userid))

	if err != nil {
		status = "error"

		return s.errorhandler.HandleFindByIdError(err.ToGRPCError(), method, "FAILED_FETCH_USER", span, &status, zap.Int("user.id", userid))
	}

	userResponse := user.Data

	s.mencache.SetCachedUserInfo(ctx, userResponse, time.Minute*5)

	logSuccess("User details fetched successfully", zap.Int("user.id", userid))

	return user.Data, nil
}
