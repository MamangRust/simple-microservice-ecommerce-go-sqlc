package orderservice

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/errorhandler"
	orderservicemapper "github.com/MamangRust/simple_microservice_ecommerce/order/internal/mapper/response/service/order"
	mencache "github.com/MamangRust/simple_microservice_ecommerce/order/internal/redis"
	orderrepository "github.com/MamangRust/simple_microservice_ecommerce/order/internal/repository/order"
	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/logger"
	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type orderQueryDeps struct {
	repository   orderrepository.OrderQueryRepository
	mapper       orderservicemapper.OrderQueryResponseMapper
	logger       logger.LoggerInterface
	errorhandler errorhandler.OrderQueryError
	mencache     mencache.OrderQueryCache
}

type orderQueryService struct {
	orderQueryRepository orderrepository.OrderQueryRepository
	mapper               orderservicemapper.OrderQueryResponseMapper
	logger               logger.LoggerInterface
	errorhandler         errorhandler.OrderQueryError
	mencache             mencache.OrderQueryCache
	observability        observability.TraceLoggerObservability
}

func NewOrderQueryService(deps *orderQueryDeps) OrderQueryService {
	observability, _ := observability.NewObservability("order-query-service", deps.logger)

	return &orderQueryService{
		orderQueryRepository: deps.repository,
		mapper:               deps.mapper,
		logger:               deps.logger,
		errorhandler:         deps.errorhandler,
		mencache:             deps.mencache,
		observability:        observability,
	}
}

func (s *orderQueryService) FindAll(ctx context.Context, req *requests.FindAllOrder) ([]*response.OrderResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	req.Page = page
	req.PageSize = pageSize

	const method = "FindAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
	)

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetOrderAllCache(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	orders, totalRecords, err := s.orderQueryRepository.FindAllOrders(ctx, req)

	if err != nil {
		status = "error"
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_ORDERS", span, &status, zap.String("error", err.Error()))
	}

	ordersResponse := s.mapper.ToOrdersResponse(orders)

	s.mencache.SetOrderAllCache(ctx, req, ordersResponse, totalRecords)

	logSuccess("Successfully retrieved orders", zap.Int("count", len(ordersResponse)), zap.Int("total_records", *totalRecords))

	return ordersResponse, totalRecords, nil
}

func (s *orderQueryService) FindByActive(ctx context.Context, req *requests.FindAllOrder) ([]*response.OrderResponseDeleteAt, *int, *response.ErrorResponse) {
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

	if data, total, found := s.mencache.GetOrderActiveCache(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	orders, totalRecords, err := s.orderQueryRepository.FindByActive(ctx, req)

	if err != nil {
		status = "error"
		errResp := response.NewErrorResponse("failed to retrieve active orders", 500)
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_ACTIVE_ORDERS", span, &status, errResp, zap.String("error", err.Error()))
	}

	ordersResponse := s.mapper.ToOrdersResponseDeleteAt(orders)

	s.mencache.SetOrderActiveCache(ctx, req, ordersResponse, totalRecords)

	logSuccess("Successfully retrieved active orders", zap.Int("count", len(ordersResponse)), zap.Int("total_records", *totalRecords))

	return ordersResponse, totalRecords, nil
}

func (s *orderQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllOrder) ([]*response.OrderResponseDeleteAt, *int, *response.ErrorResponse) {
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

	if data, total, found := s.mencache.GetOrderTrashedCache(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	orders, totalRecords, err := s.orderQueryRepository.FindByTrashed(ctx, req)

	if err != nil {
		status = "error"
		errResp := response.NewErrorResponse("failed to retrieve trashed orders", 500)
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_TRASHED_ORDERS", span, &status, errResp, zap.String("error", err.Error()))
	}

	ordersResponse := s.mapper.ToOrdersResponseDeleteAt(orders)

	s.mencache.SetOrderTrashedCache(ctx, req, ordersResponse, totalRecords)

	logSuccess("Successfully retrieved trashed orders", zap.Int("count", len(ordersResponse)), zap.Int("total_records", *totalRecords))

	return ordersResponse, totalRecords, nil
}

func (s *orderQueryService) FindById(ctx context.Context, orderID int) (*response.OrderResponse, *response.ErrorResponse) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("order_id", orderID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedOrderCache(ctx, orderID); found {
		logSuccess("Data found in cache", zap.Int("order_id", orderID))
		return data, nil
	}

	res, err := s.orderQueryRepository.FindById(ctx, orderID)

	defaultErr := response.NewErrorResponse("order not found", 404)
	if err != nil || res == nil {
		status = "error"
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_BY_ID", span, &status, defaultErr, zap.Int("order_id", orderID))
	}

	orderRes := s.mapper.ToOrderResponse(res)

	s.mencache.SetCachedOrderCache(ctx, orderRes)

	logSuccess("Successfully retrieved order", zap.Int("order_id", orderID))

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
