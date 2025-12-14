package orderitemservice

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/errorhandler"
	orderitemservicemapper "github.com/MamangRust/simple_microservice_ecommerce/order/internal/mapper/response/service/orderitem"
	mencache "github.com/MamangRust/simple_microservice_ecommerce/order/internal/redis"
	orderitemrepository "github.com/MamangRust/simple_microservice_ecommerce/order/internal/repository/order_item"
	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/logger"
	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type orderItemService struct {
	repository    orderitemrepository.OrderItemQueryRepository
	mapper        orderitemservicemapper.OrderItemResponseMapper
	logger        logger.LoggerInterface
	mencache      mencache.OrderItemQueryCache
	errorhandler  errorhandler.OrderItemQueryError
	observability observability.TraceLoggerObservability
}

type OrderItemServiceDeps struct {
	Repository   orderitemrepository.OrderItemQueryRepository
	Logger       logger.LoggerInterface
	Mencache     mencache.OrderItemQueryCache
	ErrorHandler errorhandler.OrderItemQueryError
}

func NewOrderItemService(deps *OrderItemServiceDeps) OrderItemQueryService {
	observability, _ := observability.NewObservability("order-item-query-service", deps.Logger)

	orderitemmapper := orderitemservicemapper.NewOrderItemResponseMapper()

	return &orderItemService{
		repository:    deps.Repository,
		mapper:        orderitemmapper,
		logger:        deps.Logger,
		observability: observability,
		mencache:      deps.Mencache,
		errorhandler:  deps.ErrorHandler,
	}
}

func (s *orderItemService) FindAllOrderItems(ctx context.Context, req *requests.FindAllOrderItems) ([]*response.OrderItemResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	req.Page = page
	req.PageSize = pageSize

	const method = "FindAllOrderItems"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
	)

	defer func() {
		end(status)
	}()

	// Periksa cache terlebih dahulu
	if data, total, found := s.mencache.GetCachedOrderItemsAll(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	orderItems, totalRecords, err := s.repository.FindAllOrderItems(ctx, req)

	if err != nil {
		status = "error"
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_ORDER_ITEMS", span, &status, zap.String("error", err.Error()))
	}

	orderItemsResponse := s.mapper.ToOrderItemsResponse(orderItems)

	// Isi cache setelah berhasil mengambil data
	s.mencache.SetCachedOrderItemsAll(ctx, req, orderItemsResponse, totalRecords)

	logSuccess("Successfully retrieved order items", zap.Int("count", len(orderItemsResponse)), zap.Int("total_records", *totalRecords))

	return orderItemsResponse, totalRecords, nil
}

func (s *orderItemService) FindByActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*response.OrderItemResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	req.Page = page
	req.PageSize = pageSize

	const method = "FindByActive"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
	)

	defer func() {
		end(status)
	}()

	// Periksa cache terlebih dahulu
	if data, total, found := s.mencache.GetCachedOrderItemActive(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	orderItems, totalRecords, err := s.repository.FindByActive(ctx, req)

	if err != nil {
		status = "error"
		errResp := response.NewErrorResponse("failed to retrieve active order items", 500)
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_ACTIVE_ORDER_ITEMS", span, &status, errResp, zap.String("error", err.Error()))
	}

	orderItemsResponse := s.mapper.ToOrderItemsResponseDeleteAt(orderItems)

	// Isi cache setelah berhasil mengambil data
	s.mencache.SetCachedOrderItemActive(ctx, req, orderItemsResponse, totalRecords)

	logSuccess("Successfully retrieved active order items", zap.Int("count", len(orderItemsResponse)), zap.Int("total_records", *totalRecords))

	return orderItemsResponse, totalRecords, nil
}

func (s *orderItemService) FindByTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*response.OrderItemResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	req.Page = page
	req.PageSize = pageSize

	const method = "FindByTrashed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
	)

	defer func() {
		end(status)
	}()

	// Periksa cache terlebih dahulu
	if data, total, found := s.mencache.GetCachedOrderItemTrashed(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	orderItems, totalRecords, err := s.repository.FindByTrashed(ctx, req)

	if err != nil {
		status = "error"
		errResp := response.NewErrorResponse("failed to retrieve trashed order items", 500)
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_TRASHED_ORDER_ITEMS", span, &status, errResp, zap.String("error", err.Error()))
	}

	orderItemsResponse := s.mapper.ToOrderItemsResponseDeleteAt(orderItems)

	// Isi cache setelah berhasil mengambil data
	s.mencache.SetCachedOrderItemTrashed(ctx, req, orderItemsResponse, totalRecords)

	logSuccess("Successfully retrieved trashed order items", zap.Int("count", len(orderItemsResponse)), zap.Int("total_records", *totalRecords))

	return orderItemsResponse, totalRecords, nil
}

func (s *orderItemService) FindOrderItemByOrder(ctx context.Context, orderID int) ([]*response.OrderItemResponse, *response.ErrorResponse) {
	const method = "FindOrderItemByOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("order_id", orderID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedOrderItems(ctx, orderID); found {
		logSuccess("Data found in cache", zap.Int("order_id", orderID))
		return data, nil
	}

	orderItems, err := s.repository.FindOrderItemByOrder(ctx, orderID)

	if err != nil {
		status = "error"
		errResp := response.NewErrorResponse("failed to retrieve order items for this order", 500)
		return s.errorhandler.HandleRepositoryListError(err, method, "FAILED_FIND_ORDER_ITEMS_BY_ORDER", span, &status, errResp, zap.Int("order_id", orderID))
	}

	orderItemsResponse := s.mapper.ToOrderItemsResponse(orderItems)

	s.mencache.SetCachedOrderItems(ctx, orderID, orderItemsResponse)

	logSuccess("Successfully retrieved order items for order", zap.Int("order_id", orderID), zap.Int("count", len(orderItemsResponse)))

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
