package orderservice

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/logger"
)

type OrderQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllOrder) ([]*response.OrderResponse, *int, *response.ErrorResponse)
	FindByActive(ctx context.Context, req *requests.FindAllOrder) ([]*response.OrderResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(ctx context.Context, req *requests.FindAllOrder) ([]*response.OrderResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(ctx context.Context, orderID int) (*response.OrderResponse, *response.ErrorResponse)
}

type OrderCommandService interface {
	CreateOrder(ctx context.Context, req *requests.CreateOrderRequest) (*response.OrderResponse, *response.ErrorResponse)
	UpdateOrder(ctx context.Context, req *requests.UpdateOrderRequest) (*response.OrderResponse, *response.ErrorResponse)
	TrashedOrder(ctx context.Context, orderID int) (*response.OrderResponseDeleteAt, *response.ErrorResponse)
	RestoreOrder(ctx context.Context, orderID int) (*response.OrderResponseDeleteAt, *response.ErrorResponse)
	DeleteOrderPermanent(ctx context.Context, orderID int) (bool, *response.ErrorResponse)
	RestoreAllOrder(ctx context.Context) (bool, *response.ErrorResponse)
	DeleteAllOrderPermanent(ctx context.Context) (bool, *response.ErrorResponse)
}

type OrderServiceWithLogger interface {
	OrderQueryService
	OrderCommandService
	SetLogger(logger logger.LoggerInterface)
}
