package repository

import (
	recordmapper "github.com/MamangRust/simple_microservice_ecommerce/user/internal/mapper/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/user/pkg/database/schema"
)

type Repositories interface {
	UserQueryRepo() UserQueryRepository
	UserCommandRepo() UserCommandRepository
}

type repositories struct {
	userQuery   UserQueryRepository
	userCommand UserCommandRepository
}

func (r *repositories) UserQueryRepo() UserQueryRepository {
	return r.userQuery
}

func (r *repositories) UserCommandRepo() UserCommandRepository {
	return r.userCommand
}

func NewRepositories(db *db.Queries) Repositories {
	userQueryMapper := recordmapper.NewUserQueryRecordMapper()
	userCommandMapper := recordmapper.NewUserCommandRecordMapper()

	return &repositories{
		userQuery:   NewUserQueryRepository(db, userQueryMapper),
		userCommand: NewUserCommandRepository(db, userCommandMapper),
	}
}
