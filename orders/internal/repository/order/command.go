package orderrepository

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/requests"
	orderrecordmapper "github.com/MamangRust/simple_microservice_ecommerce/order/internal/mapper/record/order"
	db "github.com/MamangRust/simple_microservice_ecommerce/order/pkg/database/schema"
	orderrepositoryerrors "github.com/MamangRust/simple_microservice_ecommerce/order/pkg/errors/order_errors/repository"
)

type orderCommandRepository struct {
	db     *db.Queries
	mapper orderrecordmapper.OrderCommandRecordMapper
}

func NewOrderCommandRepository(db *db.Queries, mapper orderrecordmapper.OrderCommandRecordMapper) OrderCommandRepository {
	return &orderCommandRepository{
		db:     db,
		mapper: mapper,
	}
}

func (r *orderCommandRepository) CreateOrder(ctx context.Context, request *requests.CreateOrderRecordRequest) (*record.OrderRecord, error) {
	req := db.CreateOrderParams{
		UserID:     int32(request.UserID),
		TotalPrice: int32(request.TotalPrice),
	}

	user, err := r.db.CreateOrder(ctx, req)

	if err != nil {
		return nil, orderrepositoryerrors.ErrCreateOrder
	}

	return r.mapper.ToOrderRecord(user), nil
}

func (r *orderCommandRepository) UpdateOrder(ctx context.Context, request *requests.UpdateOrderRecordRequest) (*record.OrderRecord, error) {
	req := db.UpdateOrderParams{
		OrderID:    int32(request.OrderID),
		TotalPrice: int32(request.TotalPrice),
	}

	res, err := r.db.UpdateOrder(ctx, req)

	if err != nil {
		return nil, orderrepositoryerrors.ErrUpdateOrder
	}

	return r.mapper.ToOrderRecord(res), nil
}

func (r *orderCommandRepository) TrashedOrder(ctx context.Context, user_id int) (*record.OrderRecord, error) {
	res, err := r.db.TrashedOrder(ctx, int32(user_id))

	if err != nil {
		return nil, orderrepositoryerrors.ErrTrashedOrder
	}

	return r.mapper.ToOrderRecord(res), nil
}

func (r *orderCommandRepository) RestoreOrder(ctx context.Context, user_id int) (*record.OrderRecord, error) {
	res, err := r.db.RestoreOrder(ctx, int32(user_id))

	if err != nil {
		return nil, orderrepositoryerrors.ErrRestoreOrder
	}

	return r.mapper.ToOrderRecord(res), nil
}

func (r *orderCommandRepository) DeleteOrderPermanent(ctx context.Context, user_id int) (bool, error) {
	err := r.db.DeleteOrderPermanently(ctx, int32(user_id))

	if err != nil {
		return false, orderrepositoryerrors.ErrDeleteOrderPermanent
	}

	return true, nil
}

func (r *orderCommandRepository) RestoreAllOrder(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllOrders(ctx)

	if err != nil {
		return false, orderrepositoryerrors.ErrRestoreAllOrder
	}
	return true, nil
}

func (r *orderCommandRepository) DeleteAllOrderPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentOrders(ctx)

	if err != nil {
		return false, orderrepositoryerrors.ErrDeleteAllOrderPermanent
	}
	return true, nil
}
