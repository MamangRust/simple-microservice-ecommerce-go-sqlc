package service

import (
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/errorhandler"
	grpcclient "github.com/MamangRust/simple_microservice_ecommerce/user/internal/grpc_client"
	responsemapper "github.com/MamangRust/simple_microservice_ecommerce/user/internal/mapper/response/service"
	mencache "github.com/MamangRust/simple_microservice_ecommerce/user/internal/redis"
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/repository"
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/hash"
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/logger"
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/observability"
)

type Service interface {
	UserQueryService
	UserCommandService
}

type service struct {
	UserQueryService
	UserCommandService
}

type GrpcClient struct {
	UserRoleClient grpcclient.UserRoleGrpcClientHandler
	RoleClient     grpcclient.RoleGrpcClientHandler
}

type Deps struct {
	Reposotories  repository.Repositories
	Hash          hash.HashPassword
	Logger        logger.LoggerInterface
	GrpcClient    GrpcClient
	Observability observability.TraceLoggerObservability
	Mencache      *mencache.Mencache
}

func NewService(deps *Deps) Service {
	mapperQuery := responsemapper.NewUserQueryResponseMapper()
	mapperCommand := responsemapper.NewUserCommandResponseMapper()
	errorhandler := errorhandler.NewErrorHandler()

	return &service{
		UserQueryService:   newUserQueryService(deps, mapperQuery, errorhandler.UserQueryError),
		UserCommandService: newUserCommandService(deps, mapperCommand, errorhandler.UserCommandError),
	}
}

func newUserQueryService(
	deps *Deps,
	mapper responsemapper.UserQueryResponseMapper,
	errorhandler errorhandler.UserQueryErrorHandler,
) UserQueryService {
	return NewUserQueryService(&userQueryDeps{
		repository:   deps.Reposotories.UserQueryRepo(),
		logger:       deps.Logger,
		mapper:       mapper,
		erorrhandler: errorhandler,
		mencache:     deps.Mencache.UserQueryCache,
	})
}

func newUserCommandService(
	deps *Deps,
	mapper responsemapper.UserCommandResponseMapper,
	errorhandler errorhandler.UserCommandErrorHandler,
) UserCommandService {
	return NewUserCommandService(&userCommandDeps{
		roleGrpcClient:        deps.GrpcClient.RoleClient,
		userRoleGrpcClient:    deps.GrpcClient.UserRoleClient,
		userQueryRepository:   deps.Reposotories.UserQueryRepo(),
		userCommandRepository: deps.Reposotories.UserCommandRepo(),
		logger:                deps.Logger,
		hash:                  deps.Hash,
		mapper:                mapper,
		errorhandler:          errorhandler,
		cacheQuery:            deps.Mencache.UserQueryCache,
		cacheCommand:          deps.Mencache.UserCommandCache,
	})
}
