package service

import (
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/errorhandler"
	grpcclient "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/grpc_client"
	mencache "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/redis"
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
	Mencache   *mencache.Mencache
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

	errorHandler := errorhandler.NewErrorHandler(deps.Logger)

	return &service{
		LoginService:         newLogin(deps, tokenService, deps.GrpcClient.UserClient, *errorHandler),
		RegistrationService:  newRegister(deps, *errorHandler),
		IdentifyService:      newIdentity(deps, tokenService, *errorHandler),
		PasswordResetService: newPasswordReset(deps, tokenService, *errorHandler),
	}
}

func newLogin(deps *Deps, tokenService *tokenService, userClient grpcclient.UserGrpcClientHandler, errorhandler errorhandler.ErrorHandler) LoginService {
	return NewLoginService(&LoginServiceDeps{
		hash:          deps.Hash,
		userClient:    userClient,
		tokenService:  tokenService,
		logger:        deps.Logger,
		errorPassword: errorhandler.PasswordError,
		errorToken:    errorhandler.TokenError,
		errorHandler:  errorhandler.LoginError,
		mencache:      deps.Mencache.LoginCache,
	})
}

func newRegister(deps *Deps, errorhandler errorhandler.ErrorHandler) RegistrationService {
	return NewRegisterService(&RegisterServiceDeps{
		userClient:        deps.GrpcClient.UserClient,
		userRoleClient:    deps.GrpcClient.UserRoleClient,
		roleClient:        deps.GrpcClient.RoleClient,
		kafka:             deps.Kafka,
		logger:            deps.Logger,
		errorRandomString: errorhandler.RandomString,
		errorMarshal:      errorhandler.MarshalError,
		errorKafka:        errorhandler.KafkaError,
		errorhandler:      errorhandler.RegisterError,
		mencache:          deps.Mencache.RegisterCache,
	})
}

func newIdentity(deps *Deps, tokenService *tokenService, errorhandler errorhandler.ErrorHandler) IdentifyService {
	return NewIdentityService(&IdentityServiceDeps{
		auth:         deps.Token,
		tokenService: tokenService,
		user:         deps.GrpcClient.UserClient,
		refreshToken: deps.Repository,
		logger:       deps.Logger,
		errorhandler: errorhandler.IdentityError,
		errorToken:   errorhandler.TokenError,
		mencache:     deps.Mencache.IdentityCache,
	})
}

func newPasswordReset(deps *Deps, tokenService *tokenService, errorhandler errorhandler.ErrorHandler) PasswordResetService {
	return NewPasswordResetService(&passwordResetServiceDeps{
		kafka:             deps.Kafka,
		UserClient:        deps.GrpcClient.UserClient,
		ResetToken:        deps.Repository,
		logger:            deps.Logger,
		errorhandler:      errorhandler.PasswordResetError,
		errorRandomString: errorhandler.RandomString,
		errorMarshal:      errorhandler.MarshalError,
		errorPassword:     errorhandler.PasswordError,
		errorKafka:        errorhandler.KafkaError,
		mencache:          deps.Mencache.PasswordResetCache,
	})
}
