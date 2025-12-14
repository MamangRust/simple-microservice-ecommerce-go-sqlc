package mencache

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/domain/response"
)

type UserCache interface {
	SetCachedUsers(ctx context.Context, req *requests.FindAllUsers, data *response.ApiResponsePaginationUser)
	SetCachedUserById(ctx context.Context, data *response.ApiResponseUser)
	SetCachedUserActive(ctx context.Context, req *requests.FindAllUsers, data *response.ApiResponsePaginationUserDeleteAt)
	SetCachedUserTrashed(ctx context.Context, req *requests.FindAllUsers, data *response.ApiResponsePaginationUserDeleteAt)

	GetCachedUsers(ctx context.Context, req *requests.FindAllUsers) (*response.ApiResponsePaginationUser, bool)
	GetCachedUserById(ctx context.Context, id int) (*response.ApiResponseUser, bool)
	GetCachedUserActive(ctx context.Context, req *requests.FindAllUsers) (*response.ApiResponsePaginationUserDeleteAt, bool)
	GetCachedUserTrashed(ctx context.Context, req *requests.FindAllUsers) (*response.ApiResponsePaginationUserDeleteAt, bool)
}

type RoleCache interface {
	SetCachedRoles(ctx context.Context, req *requests.FindAllRoles, data *response.ApiResponsePaginationRole)
	SetCachedRoleById(ctx context.Context, data *response.ApiResponseRole)
	SetCachedRoleByUserId(ctx context.Context, userId int, data *response.ApiResponsesRole)
	SetCachedRoleActive(ctx context.Context, req *requests.FindAllRoles, data *response.ApiResponsePaginationRoleDeleteAt)
	SetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRoles, data *response.ApiResponsePaginationRoleDeleteAt)

	GetCachedRoles(ctx context.Context, req *requests.FindAllRoles) (*response.ApiResponsePaginationRole, bool)
	GetCachedRoleByUserId(ctx context.Context, userId int) (*response.ApiResponsesRole, bool)
	GetCachedRoleById(ctx context.Context, id int) (*response.ApiResponseRole, bool)
	GetCachedRoleActive(ctx context.Context, req *requests.FindAllRoles) (*response.ApiResponsePaginationRoleDeleteAt, bool)
	GetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRoles) (*response.ApiResponsePaginationRoleDeleteAt, bool)
}

type ProductCache interface {
	GetCachedProducts(ctx context.Context, req *requests.FindAllProduct) (*response.ApiResponsePaginationProduct, bool)
	SetCachedProducts(ctx context.Context, req *requests.FindAllProduct, data *response.ApiResponsePaginationProduct)

	GetCachedProductActive(ctx context.Context, req *requests.FindAllProduct) (*response.ApiResponsePaginationProductDeleteAt, bool)
	SetCachedProductActive(ctx context.Context, req *requests.FindAllProduct, data *response.ApiResponsePaginationProductDeleteAt)

	GetCachedProductTrashed(ctx context.Context, req *requests.FindAllProduct) (*response.ApiResponsePaginationProductDeleteAt, bool)
	SetCachedProductTrashed(ctx context.Context, req *requests.FindAllProduct, data *response.ApiResponsePaginationProductDeleteAt)

	GetCachedProduct(ctx context.Context, productID int) (*response.ApiResponseProduct, bool)
	SetCachedProduct(ctx context.Context, data *response.ApiResponseProduct)
}

type OrderItemCache interface {
	GetCachedOrderItemsAll(ctx context.Context, req *requests.FindAllOrderItems) (*response.ApiResponsePaginationOrderItem, bool)
	SetCachedOrderItemsAll(ctx context.Context, req *requests.FindAllOrderItems, data *response.ApiResponsePaginationOrderItem)

	GetCachedOrderItemActive(ctx context.Context, req *requests.FindAllOrderItems) (*response.ApiResponsePaginationOrderItemDeleteAt, bool)
	SetCachedOrderItemActive(ctx context.Context, req *requests.FindAllOrderItems, data *response.ApiResponsePaginationOrderItemDeleteAt)

	GetCachedOrderItemTrashed(ctx context.Context, req *requests.FindAllOrderItems) (*response.ApiResponsePaginationOrderItemDeleteAt, bool)
	SetCachedOrderItemTrashed(ctx context.Context, req *requests.FindAllOrderItems, data *response.ApiResponsePaginationOrderItemDeleteAt)

	GetCachedOrderItems(ctx context.Context, orderID int) (*response.ApiResponsesOrderItem, bool)
	SetCachedOrderItems(ctx context.Context, orderID int, data *response.ApiResponsesOrderItem)
}

type OrderCache interface {
	GetOrderAllCache(ctx context.Context, req *requests.FindAllOrder) (*response.ApiResponsePaginationOrder, bool)
	SetOrderAllCache(ctx context.Context, req *requests.FindAllOrder, data *response.ApiResponsePaginationOrder)

	GetOrderActiveCache(ctx context.Context, req *requests.FindAllOrder) (*response.ApiResponsePaginationOrderDeleteAt, bool)
	SetOrderActiveCache(ctx context.Context, req *requests.FindAllOrder, data *response.ApiResponsePaginationOrderDeleteAt)

	GetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrder) (*response.ApiResponsePaginationOrderDeleteAt, bool)
	SetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrder, data *response.ApiResponsePaginationOrderDeleteAt)

	GetCachedOrderCache(ctx context.Context, orderID int) (*response.ApiResponseOrder, bool)
	SetCachedOrderCache(ctx context.Context, data *response.ApiResponseOrder)
}
