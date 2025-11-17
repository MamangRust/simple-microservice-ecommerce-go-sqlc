package orderitemservice

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
	orderitemservicemapper "github.com/MamangRust/simple_microservice_ecommerce/order/internal/mapper/response/service/orderitem"
	orderitemrepository "github.com/MamangRust/simple_microservice_ecommerce/order/internal/repository/order_item"
	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/logger"
	"go.uber.org/zap"
)

type orderItemService struct {
	repository orderitemrepository.OrderItemQueryRepository
	mapper     orderitemservicemapper.OrderItemResponseMapper
	logger     logger.LoggerInterface
}

type OrderItemServiceDeps struct {
	Repository orderitemrepository.OrderItemQueryRepository
	Logger     logger.LoggerInterface
}

func NewOrderItemService(deps *OrderItemServiceDeps) OrderItemQueryService {
	orderitemmapper := orderitemservicemapper.NewOrderItemResponseMapper()

	return &orderItemService{
		repository: deps.Repository,
		mapper:     orderitemmapper,
		logger:     deps.Logger,
	}
}

func (s *orderItemService) FindAllOrderItems(ctx context.Context, req *requests.FindAllOrderItems) ([]*response.OrderItemResponse, *int, *response.ErrorResponse) {
	s.logger.Info("Finding all order items", zap.Int("page", req.Page), zap.Int("page_size", req.PageSize))

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	req.Page = page
	req.PageSize = pageSize

	orderItems, totalRecords, err := s.repository.FindAllOrderItems(ctx, req)

	if err != nil {
		s.logger.Error("Failed to retrieve order items", zap.Int("page", page), zap.Int("page_size", pageSize), zap.Error(err))
		return nil, nil, response.NewErrorResponse("failed to retrieve order items", 500)
	}

	orderItemsResponse := s.mapper.ToOrderItemsResponse(orderItems)

	s.logger.Info("Successfully retrieved order items", zap.Int("count", len(orderItemsResponse)), zap.Int("total_records", *totalRecords))

	return orderItemsResponse, totalRecords, nil
}

func (s *orderItemService) FindByActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*response.OrderItemResponseDeleteAt, *int, *response.ErrorResponse) {
	s.logger.Info("Finding active order items", zap.Int("page", req.Page), zap.Int("page_size", req.PageSize))

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	req.Page = page
	req.PageSize = pageSize

	orderItems, totalRecords, err := s.repository.FindByActive(ctx, req)

	if err != nil {
		s.logger.Error("Failed to retrieve active order items", zap.Int("page", page), zap.Int("page_size", pageSize), zap.Error(err))
		return nil, nil, response.NewErrorResponse("failed to retrieve active order items", 500)
	}

	orderItemsResponse := s.mapper.ToOrderItemsResponseDeleteAt(orderItems)

	s.logger.Info("Successfully retrieved active order items", zap.Int("count", len(orderItemsResponse)), zap.Int("total_records", *totalRecords))

	return orderItemsResponse, totalRecords, nil
}

func (s *orderItemService) FindByTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*response.OrderItemResponseDeleteAt, *int, *response.ErrorResponse) {
	s.logger.Info("Finding trashed order items", zap.Int("page", req.Page), zap.Int("page_size", req.PageSize))

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	req.Page = page
	req.PageSize = pageSize

	orderItems, totalRecords, err := s.repository.FindByTrashed(ctx, req)

	if err != nil {
		s.logger.Error("Failed to retrieve trashed order items", zap.Int("page", page), zap.Int("page_size", pageSize), zap.Error(err))
		return nil, nil, response.NewErrorResponse("failed to retrieve trashed order items", 500)
	}

	orderItemsResponse := s.mapper.ToOrderItemsResponseDeleteAt(orderItems)

	s.logger.Info("Successfully retrieved trashed order items", zap.Int("count", len(orderItemsResponse)), zap.Int("total_records", *totalRecords))

	return orderItemsResponse, totalRecords, nil
}

func (s *orderItemService) FindOrderItemByOrder(ctx context.Context, orderID int) ([]*response.OrderItemResponse, *response.ErrorResponse) {
	s.logger.Info("Finding order items by order ID", zap.Int("order_id", orderID))

	orderItems, err := s.repository.FindOrderItemByOrder(ctx, orderID)

	if err != nil {
		s.logger.Error("Failed to retrieve order items for order", zap.Int("order_id", orderID), zap.Error(err))
		return nil, response.NewErrorResponse("failed to retrieve order items for this order", 500)
	}

	orderItemsResponse := s.mapper.ToOrderItemsResponse(orderItems)

	s.logger.Info("Successfully retrieved order items for order", zap.Int("order_id", orderID), zap.Int("count", len(orderItemsResponse)))

	return orderItemsResponse, nil
}

func (s *orderItemService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}
