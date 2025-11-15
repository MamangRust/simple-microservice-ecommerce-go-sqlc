package service

import (
	grpcclient "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/grpc_client"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/repository"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/auth"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/hash"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/kafka"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/logger"
)

type Service interface {
	LoginService
	RegistrationService
	PasswordResetService
	IdentifyService
}

type service struct {
	LoginService
	RegistrationService
	PasswordResetService
	IdentifyService
}

type Deps struct {
	Token      auth.TokenManager
	Hash       hash.HashPassword
	Repository repository.Repository
	Kafka      *kafka.Kafka
	GrpcClient GrpcClient
	Logger     logger.LoggerInterface
}

type GrpcClient struct {
	UserClient     grpcclient.UserGrpcClientHandler
	UserRoleClient grpcclient.UserRoleGrpcClientHandler
	RoleClient     grpcclient.RoleGrpcClientHandler
}

func NewService(deps *Deps) Service {
	tokenService := NewTokenService(&tokenServiceDeps{
		Token:               deps.Token,
		RefreshTokenCommand: deps.Repository,
	})

	return &service{
		LoginService:         newLogin(deps, tokenService, deps.GrpcClient.UserClient),
		RegistrationService:  newRegister(deps),
		IdentifyService:      newIdentity(deps, tokenService),
		PasswordResetService: newPasswordReset(deps, tokenService),
	}
}

func newLogin(deps *Deps, tokenService *tokenService, userClient grpcclient.UserGrpcClientHandler) LoginService {
	return NewLoginService(&LoginServiceDeps{
		hash:         deps.Hash,
		userClient:   userClient,
		tokenService: tokenService,
		logger:       deps.Logger,
	})
}

func newRegister(deps *Deps) RegistrationService {
	return NewRegisterService(&RegisterServiceDeps{
		userClient:     deps.GrpcClient.UserClient,
		userRoleClient: deps.GrpcClient.UserRoleClient,
		roleClient:     deps.GrpcClient.RoleClient,
		kafka:          deps.Kafka,
		logger:         deps.Logger,
	})
}

func newIdentity(deps *Deps, tokenService *tokenService) IdentifyService {
	return NewIdentityService(&IdentityServiceDeps{
		auth:         deps.Token,
		tokenService: tokenService,
		user:         deps.GrpcClient.UserClient,
		refreshToken: deps.Repository,
		logger:       deps.Logger,
	})
}

func newPasswordReset(deps *Deps, tokenService *tokenService) PasswordResetService {
	return NewPasswordResetService(&passwordResetServiceDeps{
		kafka:      deps.Kafka,
		UserClient: deps.GrpcClient.UserClient,
		ResetToken: deps.Repository,
		logger:     deps.Logger,
	})
}
