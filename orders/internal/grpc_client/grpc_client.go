package grpcclient

import (
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/middlewares"
	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/logger"
	pbproduct "github.com/MamangRust/simple_microservice_ecommerce_pb/product"
	pbuser "github.com/MamangRust/simple_microservice_ecommerce_pb/user"
	"google.golang.org/grpc"
)

type ServiceConnections struct {
	Product *grpc.ClientConn
	User    *grpc.ClientConn
}

type HandlerGrpcClient struct {
	ProductClient ProductGrpcClientHandler
	UserClient    UserGrpcClientHandler
}

func NewHandlerGrpcClient(service *ServiceConnections, logger logger.LoggerInterface) *HandlerGrpcClient {
	if service == nil {
		return &HandlerGrpcClient{
			ProductClient: nil,
			UserClient:    nil,
		}
	}

	errorHandling := middlewares.NewGRPCErrorHandling(logger)

	var clientQueryProduct pbproduct.ProductQueryServiceClient
	var clientCommandProduct pbproduct.ProductCommandServiceClient
	if service.Product != nil {
		clientQueryProduct = pbproduct.NewProductQueryServiceClient(service.Product)
		clientCommandProduct = pbproduct.NewProductCommandServiceClient(service.Product)
	}

	var clientUser pbuser.UserQueryServiceClient
	if service.User != nil {
		clientUser = pbuser.NewUserQueryServiceClient(service.User)
	}

	clientProductHandler := NewProductGrpcClientHandler(clientQueryProduct, clientCommandProduct, logger, errorHandling)
	clientUserHandler := NewUserGrpcClientHandler(clientUser, logger, errorHandling)

	return &HandlerGrpcClient{
		UserClient:    clientUserHandler,
		ProductClient: clientProductHandler,
	}
}
