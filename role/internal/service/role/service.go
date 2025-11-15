package roleservice

import (
	roleresponsemapper "github.com/MamangRust/simple_microservice_ecommerce/role/internal/mapper/response/role"
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
}

func NewService(repository repository.RoleRepository, logger logger.LoggerInterface) RoleService {
	mapper := roleresponsemapper.NewRoleResponseMapper()

	return &service{
		RoleQueryService:   newRoleQueryService(repository, logger, mapper),
		RoleCommandService: newRoleCommandService(repository, logger, mapper),
	}
}

func newRoleQueryService(
	repository repository.RoleRepository,
	logger logger.LoggerInterface,
	mapper roleresponsemapper.RoleResponseMapper,
) RoleQueryService {
	return NewRoleQueryService(
		&roleQueryDeps{
			repository: repository,
			logger:     logger,
			mapper:     mapper.RoleQueryResponseMapper(),
		},
	)
}

func newRoleCommandService(
	repository repository.RoleRepository,
	logger logger.LoggerInterface,
	mapper roleresponsemapper.RoleResponseMapper,
) RoleCommandService {
	return NewRoleCommandService(
		&roleCommandDeps{
			repository: repository,
			logger:     logger,
			mapper:     mapper.RoleCommandResponseMapper(),
		},
	)
}
