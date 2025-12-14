package service

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/response"
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/errorhandler"
	productresponsemapper "github.com/MamangRust/simple_microservice_ecommerce/product/internal/mapper/response"
	mencache "github.com/MamangRust/simple_microservice_ecommerce/product/internal/redis"
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/repository"
	"github.com/MamangRust/simple_microservice_ecommerce/product/pkg/logger"
	"github.com/MamangRust/simple_microservice_ecommerce/product/pkg/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type productQueryDeps struct {
	repository   repository.ProductQueryRepository
	logger       logger.LoggerInterface
	mapper       productresponsemapper.ProductQueryResponseMapper
	errorhandler errorhandler.ProductQueryError
	mencache     mencache.ProductQueryCache
}

type productQueryService struct {
	repository    repository.ProductQueryRepository
	logger        logger.LoggerInterface
	mapper        productresponsemapper.ProductQueryResponseMapper
	errorhandler  errorhandler.ProductQueryError
	mencache      mencache.ProductQueryCache
	observability observability.TraceLoggerObservability
}

func NewProductQueryService(params *productQueryDeps) ProductQueryService {
	observability, _ := observability.NewObservability("product-query-service", params.logger)

	return &productQueryService{
		repository:    params.repository,
		logger:        params.logger,
		mapper:        params.mapper,
		errorhandler:  params.errorhandler,
		mencache:      params.mencache,
		observability: observability,
	}
}

func (s *productQueryService) FindAll(ctx context.Context, req *requests.FindAllProduct) ([]*response.ProductResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	const method = "FindAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", req.Search),
	)

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedProducts(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, totalRecords, err := s.repository.FindAllProducts(ctx, req)

	if err != nil {
		status = "error"
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_PRODUCTS", span, &status, zap.String("error", err.Error()))
	}

	so := s.mapper.ToProductsResponse(res)

	s.mencache.SetCachedProducts(ctx, req, so, totalRecords)

	logSuccess("Successfully retrieved products",
		zap.Int("products_returned", len(so)),
		zap.Int("total_records", *totalRecords),
	)

	return so, totalRecords, nil
}

func (s *productQueryService) FindByID(ctx context.Context, id int) (*response.ProductResponse, *response.ErrorResponse) {
	const method = "FindByID"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("product_id", id))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedProduct(ctx, id); found {
		logSuccess("Data found in cache", zap.Int("product_id", id))
		return data, nil
	}

	res, err := s.repository.FindById(ctx, id)

	if err != nil {
		status = "error"
		defaultErr := response.NewErrorResponse("product not found", 404)
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_BY_ID", span, &status, defaultErr, zap.Int("product_id", id))
	}

	so := s.mapper.ToProductResponse(res)

	s.mencache.SetCachedProduct(ctx, so)

	logSuccess("Successfully found product by ID", zap.Int("product_id", id))
	return so, nil
}

func (s *productQueryService) FindByActive(ctx context.Context, req *requests.FindAllProduct) ([]*response.ProductResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	const method = "FindByActive"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", req.Search),
	)

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedProductActive(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, totalRecords, err := s.repository.FindByActive(ctx, req)

	if err != nil {
		status = "error"
		errResp := response.NewErrorResponse("failed to retrieve active products", 500)
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_ACTIVE_PRODUCTS", span, &status, errResp, zap.String("error", err.Error()))
	}

	so := s.mapper.ToProductsResponseDeleteAt(res)

	s.mencache.SetCachedProductActive(ctx, req, so, totalRecords)

	logSuccess("Successfully retrieved active products",
		zap.Int("products_returned", len(so)),
		zap.Int("total_records", *totalRecords),
	)

	return so, totalRecords, nil
}

func (s *productQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*response.ProductResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	search := req.Search

	const method = "FindByTrashed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedProductTrashed(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, totalRecords, err := s.repository.FindByTrashed(ctx, req)

	if err != nil {
		status = "error"
		errResp := response.NewErrorResponse("failed to retrieve trashed products", 500)
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_TRASHED_PRODUCTS", span, &status, errResp, zap.String("error", err.Error()))
	}

	so := s.mapper.ToProductsResponseDeleteAt(res)

	s.mencache.SetCachedProductTrashed(ctx, req, so, totalRecords)

	logSuccess("Successfully retrieved trashed products",
		zap.Int("products_returned", len(so)),
		zap.Int("total_records", *totalRecords),
	)

	return so, totalRecords, nil
}

func (s *productQueryService) normalizePagination(page, pageSize int) (int, int) {
	originalPage, originalPageSize := page, pageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	if page != originalPage || pageSize != originalPageSize {
		s.logger.Warn("Pagination parameters normalized",
			zap.Int("original_page", originalPage),
			zap.Int("original_page_size", originalPageSize),
			zap.Int("new_page", page),
			zap.Int("new_page_size", pageSize),
		)
	}

	return page, pageSize
}
