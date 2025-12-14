package repository

import (
	productrecordmapper "github.com/MamangRust/simple_microservice_ecommerce/product/internal/mapper/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/product/pkg/database/schema"
)

type Repositories interface {
	ProductQueryRepo() ProductQueryRepository
	ProductCommandRepo() ProductCommandRepository
}

type repositories struct {
	ProductQuery   ProductQueryRepository
	ProductCommand ProductCommandRepository
}

func (r *repositories) ProductQueryRepo() ProductQueryRepository {
	return r.ProductQuery
}

func (r *repositories) ProductCommandRepo() ProductCommandRepository {
	return r.ProductCommand
}

func NewRepositories(db *db.Queries) Repositories {
	ProductQueryMapper := productrecordmapper.NewProductQueryRecordMapper()
	ProductCommandMapper := productrecordmapper.NewProductCommandRecordMapper()

	return &repositories{
		ProductQuery:   NewProductQueryRepository(db, ProductQueryMapper),
		ProductCommand: NewProductCommandRepository(db, ProductCommandMapper),
	}
}
