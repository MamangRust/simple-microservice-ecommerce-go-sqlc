package repository

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/requests"
	productrecordmapper "github.com/MamangRust/simple_microservice_ecommerce/product/internal/mapper/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/product/pkg/database/schema"
	productCommandRepositoryerror "github.com/MamangRust/simple_microservice_ecommerce/product/pkg/errors/repository"
)

type productCommandRepository struct {
	db      *db.Queries
	mapping productrecordmapper.ProductCommandRecordMapping
}

func NewProductCommandRepository(db *db.Queries, mapping productrecordmapper.ProductCommandRecordMapping) ProductCommandRepository {
	return &productCommandRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *productCommandRepository) CreateProduct(ctx context.Context, request *requests.CreateProductRequest) (*record.ProductRecord, error) {
	req := db.CreateProductParams{
		Name:  request.Name,
		Price: int64(request.Price),
		Stock: int32(request.Stock),
	}

	product, err := r.db.CreateProduct(ctx, req)
	if err != nil {
		return nil, productCommandRepositoryerror.ErrCreateProduct
	}

	return r.mapping.ToProductRecord(product), nil
}

func (r *productCommandRepository) UpdateProduct(ctx context.Context, request *requests.UpdateProductRequest) (*record.ProductRecord, error) {
	req := db.UpdateProductParams{
		ProductID: int32(*request.ProductID),
		Name:      request.Name,
		Price:     int64(request.Price),
		Stock:     int32(request.Stock),
	}

	res, err := r.db.UpdateProduct(ctx, req)
	if err != nil {
		return nil, productCommandRepositoryerror.ErrUpdateProduct
	}

	return r.mapping.ToProductRecord(res), nil
}

func (r *productCommandRepository) UpdateProductCountStock(ctx context.Context, req *requests.UpdateProductStockRequest) (*record.ProductRecord, error) {
	res, err := r.db.UpdateProductCountStock(ctx, db.UpdateProductCountStockParams{
		ProductID: int32(req.ProductID),
		Stock:     int32(req.Stock),
	})

	if err != nil {
		return nil, productCommandRepositoryerror.ErrUpdateProductCountStock
	}

	return r.mapping.ToProductRecord(res), nil
}

func (r *productCommandRepository) TrashedProduct(ctx context.Context, product_id int) (*record.ProductRecord, error) {
	res, err := r.db.TrashProduct(ctx, int32(product_id))

	if err != nil {
		return nil, productCommandRepositoryerror.ErrTrashedProduct
	}

	return r.mapping.ToProductRecord(res), nil
}

func (r *productCommandRepository) RestoreProduct(ctx context.Context, product_id int) (*record.ProductRecord, error) {
	res, err := r.db.RestoreProduct(ctx, int32(product_id))

	if err != nil {
		return nil, productCommandRepositoryerror.ErrRestoreProduct
	}

	return r.mapping.ToProductRecord(res), nil
}

func (r *productCommandRepository) DeleteProductPermanent(ctx context.Context, product_id int) (bool, error) {
	err := r.db.DeleteProductPermanently(ctx, int32(product_id))

	if err != nil {
		return false, productCommandRepositoryerror.ErrDeleteProductPermanent
	}

	return true, nil
}

func (r *productCommandRepository) RestoreAllProducts(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllProducts(ctx)

	if err != nil {
		return false, productCommandRepositoryerror.ErrRestoreAllProducts
	}

	return true, nil
}

func (r *productCommandRepository) DeleteAllProductPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentProducts(ctx)

	if err != nil {
		return false, productCommandRepositoryerror.ErrDeleteAllProductPermanent
	}

	return true, nil
}
