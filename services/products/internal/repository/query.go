package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/requests"
	productrecordmapper "github.com/MamangRust/simple_microservice_ecommerce/product/internal/mapper/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/product/pkg/database/schema"
	productrepositoryerror "github.com/MamangRust/simple_microservice_ecommerce/product/pkg/errors/repository"
)

type productQueryRepository struct {
	db     *db.Queries
	mapper productrecordmapper.ProductQueryRecordMapping
}

func NewProductQueryRepository(db *db.Queries, mapper productrecordmapper.ProductQueryRecordMapping) ProductQueryRepository {
	return &productQueryRepository{
		db:     db,
		mapper: mapper,
	}
}

func (r *productQueryRepository) FindAllProducts(ctx context.Context, req *requests.FindAllProduct) ([]*record.ProductRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetProductsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetProducts(ctx, reqDb)
	if err != nil {
		return nil, nil, productrepositoryerror.ErrFindAllProducts
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToProductsRecordAll(res), &totalCount, nil
}

func (r *productQueryRepository) FindById(ctx context.Context, id int) (*record.ProductRecord, error) {
	res, err := r.db.GetProductByID(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, productrepositoryerror.ErrFindById
		}
		return nil, productrepositoryerror.ErrFindById
	}

	return r.mapper.ToProductRecord(res), nil
}

func (r *productQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllProduct) ([]*record.ProductRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveProductsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveProducts(ctx, reqDb)
	if err != nil {
		return nil, nil, productrepositoryerror.ErrFindByActive
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToProductsRecordActive(res), &totalCount, nil
}

func (r *productQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*record.ProductRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedProductsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedProducts(ctx, reqDb)
	if err != nil {
		return nil, nil, productrepositoryerror.ErrFindByTrashed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToProductsRecordTrashed(res), &totalCount, nil
}
