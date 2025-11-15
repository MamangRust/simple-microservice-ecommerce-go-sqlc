package service

import (
	grpcclient "github.com/MamangRust/simple_microservice_ecommerce/user/internal/grpc_client"
	responsemapper "github.com/MamangRust/simple_microservice_ecommerce/user/internal/mapper/response/service"
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/repository"
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/hash"
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/logger"
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
	Reposotories repository.Repositories
	Hash         hash.HashPassword
	Logger       logger.LoggerInterface
	GrpcClient   GrpcClient
}

func NewService(deps *Deps) Service {
	mapperQuery := responsemapper.NewUserQueryResponseMapper()
	mapperCommand := responsemapper.NewUserCommandResponseMapper()

	return &service{
		UserQueryService:   newUserQueryService(deps, mapperQuery),
		UserCommandService: newUserCommandService(deps, mapperCommand),
	}
}

func newUserQueryService(
	deps *Deps,
	mapper responsemapper.UserQueryResponseMapper,
) UserQueryService {
	return NewUserQueryService(&userQueryDeps{
		repository: deps.Reposotories.UserQueryRepo(),
		logger:     deps.Logger,
		mapper:     mapper,
	})
}

func newUserCommandService(
	deps *Deps,
	mapper responsemapper.UserCommandResponseMapper,
) UserCommandService {
	return NewUserCommandService(&userCommandDeps{
		roleGrpcClient:        deps.GrpcClient.RoleClient,
		userRoleGrpcClient:    deps.GrpcClient.UserRoleClient,
		userQueryRepository:   deps.Reposotories.UserQueryRepo(),
		userCommandRepository: deps.Reposotories.UserCommandRepo(),
		logger:                deps.Logger,
		hash:                  deps.Hash,
		mapper:                mapper,
	})
}
