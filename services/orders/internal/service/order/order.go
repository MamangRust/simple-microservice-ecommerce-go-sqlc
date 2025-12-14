package orderservice

import (
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/errorhandler"
	grpcclient "github.com/MamangRust/simple_microservice_ecommerce/order/internal/grpc_client"
	orderservicemapper "github.com/MamangRust/simple_microservice_ecommerce/order/internal/mapper/response/service/order"
	mencache "github.com/MamangRust/simple_microservice_ecommerce/order/internal/redis"
	orderrepository "github.com/MamangRust/simple_microservice_ecommerce/order/internal/repository/order"
	orderitemrepository "github.com/MamangRust/simple_microservice_ecommerce/order/internal/repository/order_item"
	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/logger"
)

type OrderService interface {
	OrderQueryService
	OrderCommandService
}

type service struct {
	OrderQueryService
	OrderCommandService
	logger logger.LoggerInterface
}

type GrpcClient struct {
	ProductClient grpcclient.ProductGrpcClientHandler
	UserClient    grpcclient.UserGrpcClientHandler
}

type Deps struct {
	OrderRepository     orderrepository.OrderRepository
	OrderItemRepository orderitemrepository.OrderItemRepository
	GrpcClient          GrpcClient
	Logger              logger.LoggerInterface
	Mencache            *mencache.Mencache
	ErrorHandler        *errorhandler.ErrorHandler
}

func NewOrderService(deps *Deps) OrderService {
	mapperQuery := orderservicemapper.NewOrderQueryResponseMapper()
	mapperCommand := orderservicemapper.NewOrderCommandResponseMapper()

	svc := &service{
		logger:              deps.Logger,
		OrderQueryService:   newOrderQueryService(deps, mapperQuery),
		OrderCommandService: newOrderCommandService(deps, mapperCommand),
	}

	return svc
}

func (s *service) SetLogger(logger logger.LoggerInterface) {
	s.logger = logger
}

func newOrderQueryService(
	deps *Deps,
	mapper orderservicemapper.OrderQueryResponseMapper,
) OrderQueryService {
	return NewOrderQueryService(&orderQueryDeps{
		repository:   deps.OrderRepository.OrderQueryRepo(),
		mapper:       mapper,
		logger:       deps.Logger,
		mencache:     deps.Mencache.OrderQueryCache,
		errorhandler: deps.ErrorHandler.OrderQueryError,
	})
}

func newOrderCommandService(deps *Deps, mapper orderservicemapper.OrderCommandResponseMapper) OrderCommandService {
	return NewOrderCommandService(&OrderCommandDeps{
		orderQueryRepository:       deps.OrderRepository.OrderQueryRepo(),
		orderCommandRepository:     deps.OrderRepository.OrderCommandRepo(),
		orderItemQueryRepository:   deps.OrderItemRepository.OrderItemQueryRepo(),
		orderItemCommandRepository: deps.OrderItemRepository.OrderItemCommandRepo(),
		mapper:                     mapper,
		userGrpcClient:             deps.GrpcClient.UserClient,
		productGrpcClient:          deps.GrpcClient.ProductClient,
		logger:                     deps.Logger,
		mencache:                   deps.Mencache.OrderCommandCache,
		errorhandler:               deps.ErrorHandler.OrderCommandError,
	})
}
