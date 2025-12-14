package service

import (
	"context"
	"time"

	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/response"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/errorhandler"
	grpcclient "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/grpc_client"
	mencache "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/redis"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/hash"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/logger"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type LoginServiceDeps struct {
	hash          hash.HashPassword
	userClient    grpcclient.UserGrpcClientHandler
	tokenService  *tokenService
	logger        logger.LoggerInterface
	errorPassword errorhandler.PasswordErrorHandler
	errorToken    errorhandler.TokenErrorHandler
	errorHandler  errorhandler.LoginErrorHandler
	mencache      mencache.LoginCache
}

type loginService struct {
	hash          hash.HashPassword
	userClient    grpcclient.UserGrpcClientHandler
	tokenService  *tokenService
	logger        logger.LoggerInterface
	errorPassword errorhandler.PasswordErrorHandler
	errorToken    errorhandler.TokenErrorHandler
	errorHandler  errorhandler.LoginErrorHandler
	observability observability.TraceLoggerObservability
	mencache      mencache.LoginCache
}

func NewLoginService(params *LoginServiceDeps) LoginService {
	observability, _ := observability.NewObservability("login-service", params.logger)

	return &loginService{
		hash:          params.hash,
		userClient:    params.userClient,
		logger:        params.logger,
		tokenService:  params.tokenService,
		errorPassword: params.errorPassword,
		errorToken:    params.errorToken,
		errorHandler:  params.errorHandler,
		observability: observability,
		mencache:      params.mencache,
	}
}

func (s *loginService) Login(ctx context.Context, request *requests.AuthRequest) (*response.TokenResponse, *response.ErrorResponse) {
	const method = "Login"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("email", request.Email))

	defer func() {
		end(status)
	}()

	if cachedToken, found := s.mencache.GetCachedLogin(ctx, request.Email); found {
		logSuccess("Successfully logged in", zap.String("email", request.Email))
		return cachedToken, nil
	}

	res, err := s.userClient.FindByEmailAndVerify(ctx, request.Email)

	if err != nil {
		status = "error"

		return s.errorHandler.HandleFindEmailError(err.ToGRPCError(), method, "LOGIN_ERR", span, &status, zap.Error(err.ToGRPCError()))
	}

	resPb := res.Data

	s.logger.Info("info password", zap.String("password", resPb.Password))

	errResp := s.hash.ComparePassword(resPb.Password, request.Password)

	if errResp != nil {
		status = "error"

		return s.errorPassword.HandleComparePasswordError(err.ToGRPCError(), method, "COMPARE_PASSWORD_ERR", span, &status, zap.Error(err.ToGRPCError()))
	}

	accessToken, errResp := s.tokenService.createAccessToken(resPb.ID)

	if errResp != nil {
		status = "error"

		return s.errorToken.HandleCreateAccessTokenError(err.ToGRPCError(), method, "CREATE_ACCESS_TOKEN_ERR", span, &status, zap.Error(err.ToGRPCError()))
	}

	refreshToken, errResp := s.tokenService.createRefreshToken(ctx, resPb.ID)

	if errResp != nil {
		status = "error"

		return s.errorToken.HandleCreateRefreshTokenError(err.ToGRPCError(), method, "CREATE_REFRESH_TOKEN_ERR", span, &status, zap.Error(err.ToGRPCError()))
	}

	tokenResp := &response.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	s.mencache.SetCachedLogin(ctx, request.Email, tokenResp, time.Minute)

	logSuccess("Successfully logged in", zap.String("email", request.Email))

	return tokenResp, nil
}
