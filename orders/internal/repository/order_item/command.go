package orderitemrepository

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/requests"
	orderitemrecordmapper "github.com/MamangRust/simple_microservice_ecommerce/order/internal/mapper/record/orderitem"
	db "github.com/MamangRust/simple_microservice_ecommerce/order/pkg/database/schema"
	orderitemrepositoryerrors "github.com/MamangRust/simple_microservice_ecommerce/order/pkg/errors/orderitem_errors/repository"
)

type orderItemCommandRepository struct {
	db     *db.Queries
	mapper orderitemrecordmapper.OrderItemCommandRecordMapper
}

func NewOrderItemCommandRepository(db *db.Queries, mapper orderitemrecordmapper.OrderItemCommandRecordMapper) OrderItemCommandRepository {
	return &orderItemCommandRepository{
		db:     db,
		mapper: mapper,
	}
}

func (r *orderItemCommandRepository) CreateOrderItem(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*record.OrderItemRecord, error) {
	res, err := r.db.CreateOrderItem(ctx, db.CreateOrderItemParams{
		OrderID:   int32(req.OrderID),
		ProductID: int32(req.ProductID),
		Quantity:  int32(req.Quantity),
		Price:     int32(req.Price),
	})

	if err != nil {
		return nil, orderitemrepositoryerrors.ErrCreateOrderItem
	}

	return r.mapper.ToOrderItemRecord(res), nil
}

func (r *orderItemCommandRepository) UpdateOrderItem(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*record.OrderItemRecord, error) {
	res, err := r.db.UpdateOrderItem(ctx, db.UpdateOrderItemParams{
		OrderItemID: int32(req.OrderItemID),
		Quantity:    int32(req.Quantity),
		Price:       int32(req.Price),
	})

	if err != nil {
		return nil, orderitemrepositoryerrors.ErrUpdateOrderItem
	}

	return r.mapper.ToOrderItemRecord(res), nil
}

func (r *orderItemCommandRepository) TrashedOrderItem(ctx context.Context, order_id int) (*record.OrderItemRecord, error) {
	res, err := r.db.TrashOrderItem(ctx, int32(order_id))

	if err != nil {
		return nil, orderitemrepositoryerrors.ErrTrashedOrderItem
	}

	return r.mapper.ToOrderItemRecord(res), nil
}

func (r *orderItemCommandRepository) RestoreOrderItem(ctx context.Context, order_id int) (*record.OrderItemRecord, error) {
	res, err := r.db.RestoreOrderItem(ctx, int32(order_id))

	if err != nil {
		return nil, orderitemrepositoryerrors.ErrRestoreOrderItem
	}

	return r.mapper.ToOrderItemRecord(res), nil
}

func (r *orderItemCommandRepository) DeleteOrderItemPermanent(ctx context.Context, order_id int) (bool, error) {
	err := r.db.DeleteOrderItemPermanently(ctx, int32(order_id))

	if err != nil {
		return false, orderitemrepositoryerrors.ErrDeleteOrderItemPermanent
	}

	return true, nil
}

func (r *orderItemCommandRepository) RestoreAllOrderItem(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllOrdersItem(ctx)

	if err != nil {
		return false, orderitemrepositoryerrors.ErrRestoreAllOrderItem
	}
	return true, nil
}

func (r *orderItemCommandRepository) DeleteAllOrderPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentOrders(ctx)

	if err != nil {
		return false, orderitemrepositoryerrors.ErrDeleteAllOrderPermanent
	}

	return true, nil
}
