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

type productCommandDeps struct {
	productQueryRepository   repository.ProductQueryRepository
	productCommandRepository repository.ProductCommandRepository
	logger                   logger.LoggerInterface
	mapper                   productresponsemapper.ProductCommandResponseMapper
	errorhandler             errorhandler.ProductCommandError
	mencache                 mencache.ProductCommandCache
}

type productCommandService struct {
	productQueryRepository   repository.ProductQueryRepository
	productCommandRepository repository.ProductCommandRepository
	logger                   logger.LoggerInterface
	mapper                   productresponsemapper.ProductCommandResponseMapper
	errorhandler             errorhandler.ProductCommandError
	mencache                 mencache.ProductCommandCache
	observability            observability.TraceLoggerObservability
}

func NewProductCommandService(params *productCommandDeps) ProductCommandService {
	observability, _ := observability.NewObservability("product-command-service", params.logger)

	return &productCommandService{
		productQueryRepository:   params.productQueryRepository,
		productCommandRepository: params.productCommandRepository,
		logger:                   params.logger,
		mapper:                   params.mapper,
		errorhandler:             params.errorhandler,
		mencache:                 params.mencache,
		observability:            observability,
	}
}

func (s *productCommandService) CreateProduct(ctx context.Context, request *requests.CreateProductRequest) (*response.ProductResponse, *response.ErrorResponse) {
	const method = "CreateProduct"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("product_name", request.Name))
	defer func() {
		end(status)
	}()

	newProduct, err := s.productCommandRepository.CreateProduct(ctx, request)
	if err != nil {
		status = "error"
		return s.errorhandler.HandleCreateProductError(err, method, "FAILED_CREATE_PRODUCT", span, &status, zap.String("product_name", request.Name))
	}

	so := s.mapper.ToProductResponse(newProduct)

	s.mencache.InvalidateAllProducts(ctx)
	s.mencache.InvalidateActiveProducts(ctx)

	logSuccess("Product created successfully",
		zap.Int("product_id", newProduct.ID),
		zap.String("product_name", newProduct.Name),
	)

	return so, nil
}

func (s *productCommandService) UpdateProduct(ctx context.Context, request *requests.UpdateProductRequest) (*response.ProductResponse, *response.ErrorResponse) {
	const method = "UpdateProduct"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("product_id", *request.ProductID))
	defer func() {
		end(status)
	}()

	_, err := s.productQueryRepository.FindById(ctx, *request.ProductID)
	if err != nil {
		status = "error"
		return s.errorhandler.HandleUpdateProductError(err, method, "FAILED_FIND_PRODUCT_FOR_UPDATE", span, &status, zap.Int("product_id", *request.ProductID))
	}

	updatedProduct, err := s.productCommandRepository.UpdateProduct(ctx, request)
	if err != nil {
		status = "error"
		return s.errorhandler.HandleUpdateProductError(err, method, "FAILED_UPDATE_PRODUCT", span, &status, zap.Int("product_id", *request.ProductID))
	}

	so := s.mapper.ToProductResponse(updatedProduct)

	s.mencache.InvalidateAllProducts(ctx)
	s.mencache.InvalidateActiveProducts(ctx)
	s.mencache.DeleteCachedProduct(ctx, *request.ProductID)

	logSuccess("Product updated successfully", zap.Int("product_id", updatedProduct.ID))

	return so, nil
}

func (s *productCommandService) UpdateProductStock(ctx context.Context, request *requests.UpdateProductStockRequest) (*response.ProductResponse, *response.ErrorResponse) {
	const method = "UpdateProductStock"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("product_id", request.ProductID))
	defer func() {
		end(status)
	}()

	// Cek apakah produk ada sebelum mencoba update stok
	_, err := s.productQueryRepository.FindById(ctx, request.ProductID)
	if err != nil {
		status = "error"
		return s.errorhandler.HandleUpdateProductError(err, method, "FAILED_FIND_PRODUCT_FOR_STOCK_UPDATE", span, &status, zap.Int("product_id", request.ProductID))
	}

	updatedProduct, err := s.productCommandRepository.UpdateProductCountStock(ctx, &requests.UpdateProductStockRequest{
		ProductID: request.ProductID,
		Stock:     request.Stock,
	})
	if err != nil {
		status = "error"
		return s.errorhandler.HandleUpdateProductError(err, method, "FAILED_UPDATE_PRODUCT_STOCK", span, &status, zap.Int("product_id", request.ProductID))
	}

	so := s.mapper.ToProductResponse(updatedProduct)

	s.mencache.InvalidateAllProducts(ctx)
	s.mencache.InvalidateActiveProducts(ctx)
	s.mencache.DeleteCachedProduct(ctx, request.ProductID)

	logSuccess("Product stock updated successfully",
		zap.Int("product_id", updatedProduct.ID),
		zap.Int("updated_stock", updatedProduct.Stock),
	)

	return so, nil
}

func (s *productCommandService) TrashedProduct(ctx context.Context, product_id int) (*response.ProductResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashedProduct"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("product_id", product_id))
	defer func() {
		end(status)
	}()

	trashedProduct, err := s.productCommandRepository.TrashedProduct(ctx, product_id)
	if err != nil {
		status = "error"
		return s.errorhandler.HandleTrashedProductError(err, method, "FAILED_TRASH_PRODUCT", span, &status, zap.Int("product_id", product_id))
	}

	so := s.mapper.ToProductResponseDeleteAt(trashedProduct)

	s.mencache.InvalidateAllProducts(ctx)
	s.mencache.InvalidateActiveProducts(ctx)
	s.mencache.DeleteCachedProduct(ctx, product_id)
	s.mencache.InvalidateTrashedProducts(ctx)

	logSuccess("Product trashed successfully", zap.Int("product_id", trashedProduct.ID))

	return so, nil
}

func (s *productCommandService) RestoreProduct(ctx context.Context, product_id int) (*response.ProductResponseDeleteAt, *response.ErrorResponse) {
	const method = "RestoreProduct"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("product_id", product_id))
	defer func() {
		end(status)
	}()

	restoredProduct, err := s.productCommandRepository.RestoreProduct(ctx, product_id)
	if err != nil {
		status = "error"
		return s.errorhandler.HandleRestoreProductError(err, method, "FAILED_RESTORE_PRODUCT", span, &status, zap.Int("product_id", product_id))
	}

	so := s.mapper.ToProductResponseDeleteAt(restoredProduct)

	s.mencache.InvalidateAllProducts(ctx)
	s.mencache.InvalidateActiveProducts(ctx)
	s.mencache.DeleteCachedProduct(ctx, product_id)
	s.mencache.InvalidateTrashedProducts(ctx)

	logSuccess("Product restored successfully", zap.Int("product_id", restoredProduct.ID))

	return so, nil
}

func (s *productCommandService) DeleteProductPermanent(ctx context.Context, product_id int) (bool, *response.ErrorResponse) {
	const method = "DeleteProductPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("product_id", product_id))
	defer func() {
		end(status)
	}()

	_, err := s.productCommandRepository.DeleteProductPermanent(ctx, product_id)
	if err != nil {
		status = "error"
		return s.errorhandler.HandleDeleteProductError(err, method, "FAILED_DELETE_PRODUCT_PERMANENT", span, &status, zap.Int("product_id", product_id))
	}

	s.mencache.InvalidateAllProducts(ctx)
	s.mencache.InvalidateActiveProducts(ctx)
	s.mencache.DeleteCachedProduct(ctx, product_id)
	s.mencache.InvalidateTrashedProducts(ctx)

	logSuccess("Product permanently deleted successfully", zap.Int("product_id", product_id))

	return true, nil
}

func (s *productCommandService) RestoreAllProduct(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "RestoreAllProduct"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	_, err := s.productCommandRepository.RestoreAllProducts(ctx)
	if err != nil {
		status = "error"
		return s.errorhandler.HandleRestoreAllProductError(err, method, "FAILED_RESTORE_ALL_PRODUCTS", span, &status)
	}

	s.mencache.InvalidateAllProducts(ctx)
	s.mencache.InvalidateActiveProducts(ctx)
	s.mencache.InvalidateTrashedProducts(ctx)

	logSuccess("All trashed products restored successfully")

	return true, nil
}

func (s *productCommandService) DeleteAllProductPermanent(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "DeleteAllProductPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	_, err := s.productCommandRepository.DeleteAllProductPermanent(ctx)
	if err != nil {
		status = "error"
		return s.errorhandler.HandleDeleteAllProductError(err, method, "FAILED_DELETE_ALL_PRODUCTS_PERMANENT", span, &status)
	}

	s.mencache.InvalidateAllProducts(ctx)
	s.mencache.InvalidateActiveProducts(ctx)
	s.mencache.InvalidateTrashedProducts(ctx)

	logSuccess("All trashed products permanently deleted successfully")

	return true, nil
}
