package handler

import (
	"strconv"

	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/middlewares"
	mencache "github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/redis"
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/pkg/logger"
	pbauth "github.com/MamangRust/simple_microservice_ecommerce_pb/auth"
	pborder "github.com/MamangRust/simple_microservice_ecommerce_pb/order"
	pborderitem "github.com/MamangRust/simple_microservice_ecommerce_pb/order_item"
	pbproduct "github.com/MamangRust/simple_microservice_ecommerce_pb/product"
	pbrole "github.com/MamangRust/simple_microservice_ecommerce_pb/role"
	pbuser "github.com/MamangRust/simple_microservice_ecommerce_pb/user"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type ServiceConnections struct {
	Auth    *grpc.ClientConn
	Role    *grpc.ClientConn
	User    *grpc.ClientConn
	Order   *grpc.ClientConn
	Product *grpc.ClientConn
}

type Deps struct {
	E                  *echo.Echo
	Logger             logger.LoggerInterface
	ServiceConnections *ServiceConnections
	Mencache           *mencache.Mencache
}

func NewHandler(deps *Deps) {
	if deps == nil || deps.ServiceConnections == nil {
		deps.Logger.Error("ServiceConnections is nil in NewHandler")
		return
	}

	grpcMiddleware := middlewares.NewGRPCErrorHandlingMiddleware(deps.Logger)

	var clientAuth pbauth.AuthServiceClient
	if deps.ServiceConnections.Auth != nil {
		clientAuth = pbauth.NewAuthServiceClient(deps.ServiceConnections.Auth)
	}

	var clientRoleQuery pbrole.RoleQueryServiceClient
	var clientRoleCommand pbrole.RoleCommandServiceClient
	if deps.ServiceConnections.Role != nil {
		clientRoleQuery = pbrole.NewRoleQueryServiceClient(deps.ServiceConnections.Role)
		clientRoleCommand = pbrole.NewRoleCommandServiceClient(deps.ServiceConnections.Role)
	}

	var clientUserQuery pbuser.UserQueryServiceClient
	var clientUserCommand pbuser.UserCommandServiceClient
	if deps.ServiceConnections.User != nil {
		clientUserQuery = pbuser.NewUserQueryServiceClient(deps.ServiceConnections.User)
		clientUserCommand = pbuser.NewUserCommandServiceClient(deps.ServiceConnections.User)
	}

	var clientOrderQuery pborder.OrderQueryServiceClient
	var clientOrderCommand pborder.OrderCommandServiceClient
	var clientOrderItem pborderitem.OrderItemServiceClient
	if deps.ServiceConnections.Order != nil {
		clientOrderQuery = pborder.NewOrderQueryServiceClient(deps.ServiceConnections.Order)
		clientOrderCommand = pborder.NewOrderCommandServiceClient(deps.ServiceConnections.Order)
		clientOrderItem = pborderitem.NewOrderItemServiceClient(deps.ServiceConnections.Order)
	}

	var clientProductQuery pbproduct.ProductQueryServiceClient
	var clientProductCommand pbproduct.ProductCommandServiceClient
	if deps.ServiceConnections.Product != nil {
		clientProductQuery = pbproduct.NewProductQueryServiceClient(deps.ServiceConnections.Product)
		clientProductCommand = pbproduct.NewProductCommandServiceClient(deps.ServiceConnections.Product)
	}

	NewAuthHandle(deps.E, clientAuth, deps.Logger, grpcMiddleware)
	NewRoleHandleApi(deps.E, clientRoleQuery, clientRoleCommand, deps.Logger, grpcMiddleware, deps.Mencache.RoleCache)
	NewUserHandleApi(deps.E, clientUserQuery, clientUserCommand, deps.Logger, grpcMiddleware, deps.Mencache.UserCache)
	NewOrderHandle(deps.E, clientOrderQuery, clientOrderCommand, deps.Logger, grpcMiddleware, deps.Mencache.OrderCache)
	NewOrderItemHandle(deps.E, clientOrderItem, deps.Logger, grpcMiddleware, deps.Mencache.OrderItemCache)
	NewProductHandle(deps.E, clientProductQuery, clientProductCommand, deps.Logger, grpcMiddleware, deps.Mencache.ProductCache)
}

func parseQueryInt(c echo.Context, key string, defaultValue int) int {
	val, err := strconv.Atoi(c.QueryParam(key))
	if err != nil || val <= 0 {
		return defaultValue
	}
	return val
}
