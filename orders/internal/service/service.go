package service

import (
	grpcclient "github.com/MamangRust/simple_microservice_ecommerce/order/internal/grpc_client"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/repository"
	orderservice "github.com/MamangRust/simple_microservice_ecommerce/order/internal/service/order"
	orderitemservice "github.com/MamangRust/simple_microservice_ecommerce/order/internal/service/orderitem"
	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/logger"
)

type Service interface {
	OrderService() orderservice.OrderService
	OrderItemService() orderitemservice.OrderItemQueryService
}

type service struct {
	orderService     orderservice.OrderService
	orderItemService orderitemservice.OrderItemQueryService
}

func (s *service) OrderService() orderservice.OrderService {
	return s.orderService
}

func (s *service) OrderItemService() orderitemservice.OrderItemQueryService {
	return s.orderItemService
}

type GrpcClient struct {
	ProductClient grpcclient.ProductGrpcClientHandler
	UserClient    grpcclient.UserGrpcClientHandler
}

type Deps struct {
	GrpcClient GrpcClient
	Repository repository.Repository
	Logger     logger.LoggerInterface
}

func NewService(deps *Deps) Service {
	return &service{
		orderService: orderservice.NewOrderService(&orderservice.Deps{
			OrderRepository:     deps.Repository,
			OrderItemRepository: deps.Repository,
			GrpcClient:          orderservice.GrpcClient(deps.GrpcClient),
			Logger:              deps.Logger,
		}),
		orderItemService: orderitemservice.NewOrderItemService(&orderitemservice.OrderItemServiceDeps{
			Repository: deps.Repository.OrderItemQueryRepo(),
			Logger:     deps.Logger,
		}),
	}
}
