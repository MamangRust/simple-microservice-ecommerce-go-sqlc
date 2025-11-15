package rolerepository

import (
	rolerecordmapper "github.com/MamangRust/simple_microservice_ecommerce/role/internal/mapper/record/role"
	db "github.com/MamangRust/simple_microservice_ecommerce/role/pkg/database/schema"
)

type RoleRepository interface {
	RoleQueryRepository
	RoleCommandRepository
}

type roleRepository struct {
	RoleQueryRepository
	RoleCommandRepository
}

func NewRoleRepository(db *db.Queries) RoleRepository {
	mapper := rolerecordmapper.NewRoleRecordMapper()

	return &roleRepository{
		RoleQueryRepository:   NewRoleQueryRepository(db, mapper.RoleQueryRecordMapper()),
		RoleCommandRepository: NewRoleCommandRepository(db, mapper.RoleCommandRecordMapper()),
	}
}
