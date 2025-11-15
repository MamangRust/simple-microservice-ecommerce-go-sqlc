package service

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/response"
	productresponsemapper "github.com/MamangRust/simple_microservice_ecommerce/product/internal/mapper/response"
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/repository"
	"github.com/MamangRust/simple_microservice_ecommerce/product/pkg/logger"
	"go.uber.org/zap"
)

type productQueryDeps struct {
	repository repository.ProductQueryRepository
	logger     logger.LoggerInterface
	mapper     productresponsemapper.ProductQueryResponseMapper
}

type productQueryService struct {
	repository repository.ProductQueryRepository
	logger     logger.LoggerInterface
	mapper     productresponsemapper.ProductQueryResponseMapper
}

func NewProductQueryService(params *productQueryDeps) ProductQueryService {
	return &productQueryService{
		repository: params.repository,
		logger:     params.logger,
		mapper:     params.mapper,
	}
}

func (s *productQueryService) FindAll(ctx context.Context, req *requests.FindAllProduct) ([]*response.ProductResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	req.Page = page
	req.PageSize = pageSize

	res, totalRecords, err := s.repository.FindAllProducts(ctx, req)

	if err != nil {
		s.logger.Error("Failed to retrieve all products from repository", zap.Error(err))
		return nil, nil, response.NewErrorResponse("failed to retrieve products", 500)
	}

	so := s.mapper.ToProductsResponse(res)

	s.logger.Info("Successfully retrieved products",
		zap.Int("products_returned", len(so)),
		zap.Int("total_records", *totalRecords),
	)

	return so, totalRecords, nil
}

func (s *productQueryService) FindByID(ctx context.Context, id int) (*response.ProductResponse, *response.ErrorResponse) {
	res, err := s.repository.FindById(ctx, id)

	if err != nil {
		s.logger.Warn("Product not found by ID", zap.Int("product_id", id), zap.Error(err))
		return nil, response.NewErrorResponse("product not found", 404)
	}

	so := s.mapper.ToProductResponse(res)

	s.logger.Info("Successfully found product by ID", zap.Int("product_id", id))
	return so, nil
}

func (s *productQueryService) FindByActive(ctx context.Context, req *requests.FindAllProduct) ([]*response.ProductResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	req.Page = page
	req.PageSize = pageSize

	res, totalRecords, err := s.repository.FindByActive(ctx, req)

	if err != nil {
		s.logger.Error("Failed to retrieve active products from repository", zap.Error(err))
		return nil, nil, response.NewErrorResponse("failed to retrieve active products", 500)
	}

	so := s.mapper.ToProductsResponseDeleteAt(res)

	s.logger.Info("Successfully retrieved active products",
		zap.Int("products_returned", len(so)),
		zap.Int("total_records", *totalRecords),
	)

	return so, totalRecords, nil
}

func (s *productQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*response.ProductResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	req.Page = page
	req.PageSize = pageSize

	res, totalRecords, err := s.repository.FindByTrashed(ctx, req)

	if err != nil {
		s.logger.Error("Failed to retrieve trashed products from repository", zap.Error(err))
		return nil, nil, response.NewErrorResponse("failed to retrieve trashed products", 500)
	}

	so := s.mapper.ToProductsResponseDeleteAt(res)

	s.logger.Info("Successfully retrieved trashed products",
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
