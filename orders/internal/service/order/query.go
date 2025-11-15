package orderservice

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
	orderservicemapper "github.com/MamangRust/simple_microservice_ecommerce/order/internal/mapper/response/service/order"
	orderrepository "github.com/MamangRust/simple_microservice_ecommerce/order/internal/repository/order"
	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/logger"
	"go.uber.org/zap"
)

type orderQueryDeps struct {
	repository orderrepository.OrderQueryRepository
	mapper     orderservicemapper.OrderQueryResponseMapper
	logger     logger.LoggerInterface
}

type orderQueryService struct {
	orderQueryRepository orderrepository.OrderQueryRepository
	mapper               orderservicemapper.OrderQueryResponseMapper
	logger               logger.LoggerInterface
}

func NewOrderQueryService(deps *orderQueryDeps) OrderQueryService {
	return &orderQueryService{
		orderQueryRepository: deps.repository,
		mapper:               deps.mapper,
		logger:               deps.logger,
	}
}

func (s *orderQueryService) FindAll(ctx context.Context, req *requests.FindAllOrder) ([]*response.OrderResponse, *int, *response.ErrorResponse) {
	s.logger.Info("Finding all orders", zap.Int("page", req.Page), zap.Int("page_size", req.PageSize))

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	req.Page = page
	req.PageSize = pageSize

	orders, totalRecords, err := s.orderQueryRepository.FindAllOrders(ctx, req)

	if err != nil {
		s.logger.Error("Failed to retrieve orders", zap.Int("page", page), zap.Int("page_size", pageSize), zap.Error(err))
		return nil, nil, response.NewErrorResponse("failed to retrieve orders", 500)
	}

	ordersResponse := s.mapper.ToOrdersResponse(orders)

	s.logger.Info("Successfully retrieved orders", zap.Int("count", len(ordersResponse)), zap.Int("total_records", *totalRecords))

	return ordersResponse, totalRecords, nil
}

func (s *orderQueryService) FindByActive(ctx context.Context, req *requests.FindAllOrder) ([]*response.OrderResponseDeleteAt, *int, *response.ErrorResponse) {
	s.logger.Info("Finding active orders", zap.Int("page", req.Page), zap.Int("page_size", req.PageSize))

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	req.Page = page
	req.PageSize = pageSize

	orders, totalRecords, err := s.orderQueryRepository.FindByActive(ctx, req)

	if err != nil {
		s.logger.Error("Failed to retrieve active orders", zap.Int("page", page), zap.Int("page_size", pageSize), zap.Error(err))
		return nil, nil, response.NewErrorResponse("failed to retrieve orders", 500)
	}

	ordersResponse := s.mapper.ToOrdersResponseDeleteAt(orders)

	s.logger.Info("Successfully retrieved active orders", zap.Int("count", len(ordersResponse)), zap.Int("total_records", *totalRecords))

	return ordersResponse, totalRecords, nil
}

func (s *orderQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllOrder) ([]*response.OrderResponseDeleteAt, *int, *response.ErrorResponse) {
	s.logger.Info("Finding trashed orders", zap.Int("page", req.Page), zap.Int("page_size", req.PageSize))

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	req.Page = page
	req.PageSize = pageSize

	orders, totalRecords, err := s.orderQueryRepository.FindByTrashed(ctx, req)

	if err != nil {
		s.logger.Error("Failed to retrieve trashed orders", zap.Int("page", page), zap.Int("page_size", pageSize), zap.Error(err))
		return nil, nil, response.NewErrorResponse("failed to retrieve orders", 500)
	}

	ordersResponse := s.mapper.ToOrdersResponseDeleteAt(orders)

	s.logger.Info("Successfully retrieved trashed orders", zap.Int("count", len(ordersResponse)), zap.Int("total_records", *totalRecords))

	return ordersResponse, totalRecords, nil
}

func (s *orderQueryService) FindById(ctx context.Context, orderID int) (*response.OrderResponse, *response.ErrorResponse) {
	s.logger.Info("Finding order by ID", zap.Int("order_id", orderID))

	res, err := s.orderQueryRepository.FindById(ctx, orderID)

	if err != nil {
		s.logger.Error("Order not found", zap.Int("order_id", orderID), zap.Error(err))
		return nil, response.NewErrorResponse("order not found", 500)
	}

	orderRes := s.mapper.ToOrderResponse(res)

	s.logger.Info("Successfully retrieved order", zap.Int("order_id", orderID))

	return orderRes, nil
}

func (s *orderQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}
