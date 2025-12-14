package orderitemrepository

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/requests"
)

type OrderItemQueryRepository interface {
	FindOrderItemByOrder(ctx context.Context, orderID int) ([]*record.OrderItemRecord, error)
	CalculateTotalPrice(ctx context.Context, orderID int) (*int32, error)
	FindAllOrderItems(ctx context.Context, req *requests.FindAllOrderItems) ([]*record.OrderItemRecord, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*record.OrderItemRecord, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*record.OrderItemRecord, *int, error)
}

type OrderItemCommandRepository interface {
	CreateOrderItem(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*record.OrderItemRecord, error)
	UpdateOrderItem(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*record.OrderItemRecord, error)
	TrashedOrderItem(ctx context.Context, orderID int) (*record.OrderItemRecord, error)
	RestoreOrderItem(ctx context.Context, orderID int) (*record.OrderItemRecord, error)
	DeleteOrderItemPermanent(ctx context.Context, orderID int) (bool, error)
	RestoreAllOrderItem(ctx context.Context) (bool, error)
	DeleteAllOrderPermanent(ctx context.Context) (bool, error)
}
