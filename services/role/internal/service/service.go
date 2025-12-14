package service

import (
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/errorhandler"
	mencache "github.com/MamangRust/simple_microservice_ecommerce/role/internal/redis"
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/repository"
	roleservice "github.com/MamangRust/simple_microservice_ecommerce/role/internal/service/role"
	userroleservice "github.com/MamangRust/simple_microservice_ecommerce/role/internal/service/user_role"
	"github.com/MamangRust/simple_microservice_ecommerce/role/pkg/logger"
)

type Service interface {
	userroleservice.UserRoleService
	roleservice.RoleService
}

type service struct {
	userroleservice.UserRoleService
	roleservice.RoleService
}

type Deps struct {
	Repository repository.Repository
	Logger     logger.LoggerInterface
	Mencache   *mencache.Mencache
}

func NewService(deps *Deps) Service {
	errorhandler := errorhandler.NewErrorHandler(deps.Logger)

	return &service{
		UserRoleService: userroleservice.NewUserRoleService(&userroleservice.UserRoleServiceDeps{
			Logger:       deps.Logger,
			Repository:   deps.Repository,
			Errorhandler: errorhandler.UserRoleCommandError,
		}),
		RoleService: roleservice.NewService(&roleservice.Deps{
			RoleRepository: deps.Repository,
			Logger:         deps.Logger,
			ErrorHandler:   errorhandler,
			Mencache:       deps.Mencache,
		}),
	}
}
