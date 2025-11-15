package orderrepository

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/record"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/requests"
)

type OrderQueryRepository interface {
	FindAllOrders(ctx context.Context, req *requests.FindAllOrder) ([]*record.OrderRecord, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllOrder) ([]*record.OrderRecord, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllOrder) ([]*record.OrderRecord, *int, error)
	FindById(ctx context.Context, orderID int) (*record.OrderRecord, error)
}

type OrderCommandRepository interface {
	CreateOrder(ctx context.Context, req *requests.CreateOrderRecordRequest) (*record.OrderRecord, error)
	UpdateOrder(ctx context.Context, req *requests.UpdateOrderRecordRequest) (*record.OrderRecord, error)
	TrashedOrder(ctx context.Context, orderID int) (*record.OrderRecord, error)
	RestoreOrder(ctx context.Context, orderID int) (*record.OrderRecord, error)
	DeleteOrderPermanent(ctx context.Context, orderID int) (bool, error)
	RestoreAllOrder(ctx context.Context) (bool, error)
	DeleteAllOrderPermanent(ctx context.Context) (bool, error)
}
