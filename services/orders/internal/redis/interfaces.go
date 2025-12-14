package mencache

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
)

type OrderItemQueryCache interface {
	GetCachedOrderItemsAll(ctx context.Context, req *requests.FindAllOrderItems) ([]*response.OrderItemResponse, *int, bool)
	SetCachedOrderItemsAll(ctx context.Context, req *requests.FindAllOrderItems, data []*response.OrderItemResponse, total *int)

	GetCachedOrderItemActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*response.OrderItemResponseDeleteAt, *int, bool)
	SetCachedOrderItemActive(ctx context.Context, req *requests.FindAllOrderItems, data []*response.OrderItemResponseDeleteAt, total *int)

	GetCachedOrderItemTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*response.OrderItemResponseDeleteAt, *int, bool)
	SetCachedOrderItemTrashed(ctx context.Context, req *requests.FindAllOrderItems, data []*response.OrderItemResponseDeleteAt, total *int)

	GetCachedOrderItems(ctx context.Context, orderID int) ([]*response.OrderItemResponse, bool)
	SetCachedOrderItems(ctx context.Context, orderID int, data []*response.OrderItemResponse)
}

type OrderQueryCache interface {
	GetOrderAllCache(ctx context.Context, req *requests.FindAllOrder) ([]*response.OrderResponse, *int, bool)
	SetOrderAllCache(ctx context.Context, req *requests.FindAllOrder, data []*response.OrderResponse, total *int)

	GetOrderActiveCache(ctx context.Context, req *requests.FindAllOrder) ([]*response.OrderResponseDeleteAt, *int, bool)
	SetOrderActiveCache(ctx context.Context, req *requests.FindAllOrder, data []*response.OrderResponseDeleteAt, total *int)

	GetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrder) ([]*response.OrderResponseDeleteAt, *int, bool)
	SetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrder, data []*response.OrderResponseDeleteAt, total *int)

	GetCachedOrderCache(ctx context.Context, orderID int) (*response.OrderResponse, bool)
	SetCachedOrderCache(ctx context.Context, data *response.OrderResponse)
}

type OrderCommandCache interface {
	DeleteOrderCache(ctx context.Context, orderID int)

	InvalidateAllOrders(ctx context.Context)
	InvalidateActiveOrders(ctx context.Context)
	InvalidateTrashedOrders(ctx context.Context)
}
