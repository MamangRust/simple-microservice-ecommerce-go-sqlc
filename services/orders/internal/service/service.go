package service

import (
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/errorhandler"
	grpcclient "github.com/MamangRust/simple_microservice_ecommerce/order/internal/grpc_client"
	mencache "github.com/MamangRust/simple_microservice_ecommerce/order/internal/redis"
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
	Mencache   *mencache.Mencache
}

func NewService(deps *Deps) Service {
	errorhandler := errorhandler.NewErrorHandler(deps.Logger)

	return &service{
		orderService: orderservice.NewOrderService(&orderservice.Deps{
			OrderRepository:     deps.Repository,
			OrderItemRepository: deps.Repository,
			GrpcClient:          orderservice.GrpcClient(deps.GrpcClient),
			Logger:              deps.Logger,
			Mencache:            deps.Mencache,
			ErrorHandler:        errorhandler,
		}),
		orderItemService: orderitemservice.NewOrderItemService(&orderitemservice.OrderItemServiceDeps{
			Repository:   deps.Repository.OrderItemQueryRepo(),
			Logger:       deps.Logger,
			ErrorHandler: errorhandler.OrderItemQueryError,
			Mencache:     deps.Mencache.OrderItemQueryCache,
		}),
	}
}
