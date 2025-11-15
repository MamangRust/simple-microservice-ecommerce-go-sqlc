package orderitemrepository

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/requests"
	orderitemrecordmapper "github.com/MamangRust/simple_microservice_ecommerce/order/internal/mapper/record/orderitem"
	db "github.com/MamangRust/simple_microservice_ecommerce/order/pkg/database/schema"
	orderitemrepositoryerrors "github.com/MamangRust/simple_microservice_ecommerce/order/pkg/errors/orderitem_errors/repository"
)

type orderItemQueryRepository struct {
	db     *db.Queries
	mapper orderitemrecordmapper.OrderItemQueryRecordMapper
}

func NewOrderItemQueryRepository(db *db.Queries, mapper orderitemrecordmapper.OrderItemQueryRecordMapper) OrderItemQueryRepository {
	return &orderItemQueryRepository{
		db:     db,
		mapper: mapper,
	}
}

func (r *orderItemQueryRepository) FindAllOrderItems(ctx context.Context, req *requests.FindAllOrderItems) ([]*record.OrderItemRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrderItemsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrderItems(ctx, reqDb)

	if err != nil {
		return nil, nil, orderitemrepositoryerrors.ErrFindAllOrderItems
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToOrderItemsRecordPagination(res), &totalCount, nil
}

func (r *orderItemQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*record.OrderItemRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrderItemsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrderItemsActive(ctx, reqDb)

	if err != nil {
		return nil, nil, orderitemrepositoryerrors.ErrFindByActive
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToOrderItemsRecordActivePagination(res), &totalCount, nil
}

func (r *orderItemQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*record.OrderItemRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrderItemsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrderItemsTrashed(ctx, reqDb)

	if err != nil {
		return nil, nil, orderitemrepositoryerrors.ErrFindByTrashed
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToOrderItemsRecordTrashedPagination(res), &totalCount, nil
}

func (r *orderItemQueryRepository) CalculateTotalPrice(ctx context.Context, order_id int) (*int32, error) {
	res, err := r.db.CalculateTotalPrice(ctx, int32(order_id))

	if err != nil {
		return nil, orderitemrepositoryerrors.ErrCalculateTotalPrice
	}

	return &res, nil

}

func (r *orderItemQueryRepository) FindOrderItemByOrder(ctx context.Context, order_id int) ([]*record.OrderItemRecord, error) {
	res, err := r.db.GetOrderItemsByOrder(ctx, int32(order_id))

	if err != nil {
		return nil, orderitemrepositoryerrors.ErrFindOrderItemByOrder
	}

	return r.mapper.ToOrderItemsRecord(res), nil
}
