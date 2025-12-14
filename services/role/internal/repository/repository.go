package repository

import (
	userrolerecordmapper "github.com/MamangRust/simple_microservice_ecommerce/role/internal/mapper/record/user_role"
	rolerepository "github.com/MamangRust/simple_microservice_ecommerce/role/internal/repository/role"
	userrolerepository "github.com/MamangRust/simple_microservice_ecommerce/role/internal/repository/user_role"
	db "github.com/MamangRust/simple_microservice_ecommerce/role/pkg/database/schema"
)

type Repository interface {
	rolerepository.RoleRepository
	userrolerepository.UserRoleRepository
}

type repositories struct {
	rolerepository.RoleRepository
	userrolerepository.UserRoleRepository
}

func NewRepositories(db *db.Queries) Repository {

	mapperUserRole := userrolerecordmapper.NewUserRoleRecordMapper()
	return &repositories{
		RoleRepository:     rolerepository.NewRoleRepository(db),
		UserRoleRepository: userrolerepository.NewUserRoleRepository(db, mapperUserRole),
	}
}
