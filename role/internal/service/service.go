package service

import (
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
}

func NewService(repository repository.Repository, logger logger.LoggerInterface) Service {
	return &service{
		UserRoleService: userroleservice.NewUserRoleService(repository, logger),
		RoleService:     roleservice.NewService(repository, logger),
	}
}
