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

type productCommandDeps struct {
	productQueryRepository   repository.ProductQueryRepository
	productCommandRepository repository.ProductCommandRepository
	logger                   logger.LoggerInterface
	mapper                   productresponsemapper.ProductCommandResponseMapper
}

type productCommandService struct {
	productQueryRepository   repository.ProductQueryRepository
	productCommandRepository repository.ProductCommandRepository
	logger                   logger.LoggerInterface
	mapper                   productresponsemapper.ProductCommandResponseMapper
}

func NewProductCommandService(params *productCommandDeps) ProductCommandService {
	return &productCommandService{
		productQueryRepository:   params.productQueryRepository,
		productCommandRepository: params.productCommandRepository,
		logger:                   params.logger,
		mapper:                   params.mapper,
	}
}

func (s *productCommandService) CreateProduct(ctx context.Context, request *requests.CreateProductRequest) (*response.ProductResponse, *response.ErrorResponse) {
	newProduct, err := s.productCommandRepository.CreateProduct(ctx, request)
	if err != nil {
		s.logger.Error("Failed to create product in repository",
			zap.Error(err),
			zap.String("product_name", request.Name),
		)
		return nil, response.NewErrorResponse("failed to create product", 500)
	}

	s.logger.Info("Product created successfully",
		zap.Int("product_id", newProduct.ID),
		zap.String("product_name", newProduct.Name),
	)

	so := s.mapper.ToProductResponse(newProduct)
	return so, nil
}

func (s *productCommandService) UpdateProduct(ctx context.Context, request *requests.UpdateProductRequest) (*response.ProductResponse, *response.ErrorResponse) {
	_, err := s.productQueryRepository.FindById(ctx, *request.ProductID)
	if err != nil {
		s.logger.Warn("Product not found for update", zap.Error(err), zap.Int("product_id", *request.ProductID))
		return nil, response.NewErrorResponse("product not found", 404)
	}

	updatedProduct, err := s.productCommandRepository.UpdateProduct(ctx, request)
	if err != nil {
		s.logger.Error("Failed to update product in repository",
			zap.Error(err),
			zap.Int("product_id", *request.ProductID),
		)
		return nil, response.NewErrorResponse("failed to update product", 500)
	}

	s.logger.Info("Product updated successfully", zap.Int("product_id", updatedProduct.ID))

	so := s.mapper.ToProductResponse(updatedProduct)
	return so, nil
}

func (s *productCommandService) UpdateProductStock(ctx context.Context, request *requests.UpdateProductStockRequest) (*response.ProductResponse, *response.ErrorResponse) {
	_, err := s.productQueryRepository.FindById(ctx, request.ProductID)
	if err != nil {
		s.logger.Warn("Product not found for stock update", zap.Error(err), zap.Int("product_id", request.ProductID))
		return nil, response.NewErrorResponse("product not found", 404)
	}

	updatedProduct, err := s.productCommandRepository.UpdateProductCountStock(ctx, &requests.UpdateProductStockRequest{
		ProductID: request.ProductID,
		Stock:     request.Stock,
	})
	if err != nil {
		s.logger.Error("Failed to update product stock in repository",
			zap.Error(err),
			zap.Int("product_id", request.ProductID),
		)
		return nil, response.NewErrorResponse("failed to update product stock", 500)
	}

	s.logger.Info("Product stock updated successfully",
		zap.Int("product_id", updatedProduct.ID),
		zap.Int("updated_stock", updatedProduct.Stock),
	)

	so := s.mapper.ToProductResponse(updatedProduct)
	return so, nil
}

func (s *productCommandService) TrashedProduct(ctx context.Context, product_id int) (*response.ProductResponseDeleteAt, *response.ErrorResponse) {
	trashedProduct, err := s.productCommandRepository.TrashedProduct(ctx, product_id)
	if err != nil {
		s.logger.Error("Failed to trash product in repository",
			zap.Error(err),
			zap.Int("product_id", product_id),
		)
		return nil, response.NewErrorResponse("failed to trash product", 500)
	}

	s.logger.Info("Product trashed successfully", zap.Int("product_id", trashedProduct.ID))

	so := s.mapper.ToProductResponseDeleteAt(trashedProduct)
	return so, nil
}

func (s *productCommandService) RestoreProduct(ctx context.Context, product_id int) (*response.ProductResponseDeleteAt, *response.ErrorResponse) {
	restoredProduct, err := s.productCommandRepository.RestoreProduct(ctx, product_id)
	if err != nil {
		s.logger.Error("Failed to restore product in repository",
			zap.Error(err),
			zap.Int("product_id", product_id),
		)
		return nil, response.NewErrorResponse("failed to restore product", 500)
	}

	s.logger.Info("Product restored successfully", zap.Int("product_id", restoredProduct.ID))

	so := s.mapper.ToProductResponseDeleteAt(restoredProduct)
	return so, nil
}

func (s *productCommandService) DeleteProductPermanent(ctx context.Context, product_id int) (bool, *response.ErrorResponse) {
	_, err := s.productCommandRepository.DeleteProductPermanent(ctx, product_id)
	if err != nil {
		s.logger.Error("Failed to permanently delete product in repository",
			zap.Error(err),
			zap.Int("product_id", product_id),
		)
		return false, response.NewErrorResponse("failed to delete product permanently", 500)
	}

	s.logger.Info("Product permanently deleted successfully", zap.Int("product_id", product_id))
	return true, nil
}

func (s *productCommandService) RestoreAllProduct(ctx context.Context) (bool, *response.ErrorResponse) {
	_, err := s.productCommandRepository.RestoreAllProducts(ctx)

	if err != nil {
		s.logger.Error("Failed to restore all products in repository", zap.Error(err))
		return false, response.NewErrorResponse("failed to restore all products", 500)
	}

	s.logger.Info("All trashed products restored successfully")
	return true, nil
}

func (s *productCommandService) DeleteAllProductPermanent(ctx context.Context) (bool, *response.ErrorResponse) {
	_, err := s.productCommandRepository.DeleteAllProductPermanent(ctx)
	if err != nil {
		s.logger.Error("Failed to permanently delete all products in repository", zap.Error(err))
		return false, response.NewErrorResponse("failed to delete all products permanently", 500)
	}

	s.logger.Info("All trashed products permanently deleted successfully")
	return true, nil
}
