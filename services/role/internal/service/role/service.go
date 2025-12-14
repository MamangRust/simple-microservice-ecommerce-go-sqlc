package roleservice

import (
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/errorhandler"
	roleresponsemapper "github.com/MamangRust/simple_microservice_ecommerce/role/internal/mapper/response/role"
	mencache "github.com/MamangRust/simple_microservice_ecommerce/role/internal/redis"
	repository "github.com/MamangRust/simple_microservice_ecommerce/role/internal/repository/role"
	"github.com/MamangRust/simple_microservice_ecommerce/role/pkg/logger"
)

type RoleService interface {
	RoleQueryService
	RoleCommandService
}

type service struct {
	RoleQueryService
	RoleCommandService
}

type Deps struct {
	RoleRepository repository.RoleRepository
	Logger         logger.LoggerInterface
	ErrorHandler   *errorhandler.ErrorHandler
	Mencache       *mencache.Mencache
}

func NewService(deps *Deps) RoleService {
	mapper := roleresponsemapper.NewRoleResponseMapper()

	return &service{
		RoleQueryService:   newRoleQueryService(deps.RoleRepository, deps.Logger, mapper, deps.ErrorHandler.RoleQueryError, deps.Mencache.RoleQueryCache),
		RoleCommandService: newRoleCommandService(deps.RoleRepository, deps.Logger, mapper, deps.ErrorHandler.RoleCommandError, deps.Mencache),
	}
}

func newRoleQueryService(
	repository repository.RoleRepository,
	logger logger.LoggerInterface,
	mapper roleresponsemapper.RoleResponseMapper,
	errorhandler errorhandler.RoleQueryErrorHandler,
	mencache mencache.RoleQueryCache,
) RoleQueryService {
	return NewRoleQueryService(
		&roleQueryDeps{
			repository:   repository,
			logger:       logger,
			mapper:       mapper.RoleQueryResponseMapper(),
			errorhandler: errorhandler,
			mencache:     mencache,
		},
	)
}

func newRoleCommandService(
	repository repository.RoleRepository,
	logger logger.LoggerInterface,
	mapper roleresponsemapper.RoleResponseMapper,
	errorhandler errorhandler.RoleCommandErrorHandler,
	mencache *mencache.Mencache,
) RoleCommandService {
	return NewRoleCommandService(
		&roleCommandDeps{
			repository:   repository,
			logger:       logger,
			mapper:       mapper.RoleCommandResponseMapper(),
			errorhandler: errorhandler,
			cacheQuery:   mencache.RoleQueryCache,
			cacheCommand: mencache.RoleCommandCache,
		},
	)
}
