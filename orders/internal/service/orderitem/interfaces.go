package orderitemservice

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
	orderitemrepository "github.com/MamangRust/simple_microservice_ecommerce/order/internal/repository/order_item"
	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/logger"
)

type OrderItemQueryService interface {
	FindAllOrderItems(ctx context.Context, req *requests.FindAllOrderItems) ([]*response.OrderItemResponse, *int, *response.ErrorResponse)
	FindByActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*response.OrderItemResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*response.OrderItemResponseDeleteAt, *int, *response.ErrorResponse)
	FindOrderItemByOrder(ctx context.Context, orderID int) ([]*response.OrderItemResponse, *response.ErrorResponse)
}

type OrderItemServiceDeps struct {
	Repository orderitemrepository.OrderItemQueryRepository
	Logger     logger.LoggerInterface
}
