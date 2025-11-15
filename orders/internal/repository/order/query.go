package orderrepository

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/requests"
	orderrecordmapper "github.com/MamangRust/simple_microservice_ecommerce/order/internal/mapper/record/order"
	db "github.com/MamangRust/simple_microservice_ecommerce/order/pkg/database/schema"
	orderrepositoryerrors "github.com/MamangRust/simple_microservice_ecommerce/order/pkg/errors/order_errors/repository"
)

type orderQueryRepository struct {
	db     *db.Queries
	mapper orderrecordmapper.OrderQueryRecordMapper
}

func NewOrderQueryRepository(db *db.Queries, mapper orderrecordmapper.OrderQueryRecordMapper) OrderQueryRepository {
	return &orderQueryRepository{
		db:     db,
		mapper: mapper,
	}
}

func (r *orderQueryRepository) FindAllOrders(ctx context.Context, req *requests.FindAllOrder) ([]*record.OrderRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrdersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrders(ctx, reqDb)

	if err != nil {
		return nil, nil, orderrepositoryerrors.ErrFindAllOrders
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToOrdersRecordPagination(res), &totalCount, nil
}

func (r *orderQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllOrder) ([]*record.OrderRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrdersActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrdersActive(ctx, reqDb)

	if err != nil {
		return nil, nil, orderrepositoryerrors.ErrFindByActive
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToOrdersRecordActivePagination(res), &totalCount, nil
}

func (r *orderQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllOrder) ([]*record.OrderRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrdersTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrdersTrashed(ctx, reqDb)

	if err != nil {
		return nil, nil, orderrepositoryerrors.ErrFindByTrashed
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToOrdersRecordTrashedPagination(res), &totalCount, nil
}

func (r *orderQueryRepository) FindById(ctx context.Context, user_id int) (*record.OrderRecord, error) {
	res, err := r.db.GetOrderByID(ctx, int32(user_id))

	if err != nil {
		return nil, orderrepositoryerrors.ErrFindById
	}

	return r.mapper.ToOrderRecord(res), nil
}
