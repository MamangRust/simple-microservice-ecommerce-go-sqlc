package handler

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/requests"
	productprotomapper "github.com/MamangRust/simple_microservice_ecommerce/product/internal/mapper/proto/product"
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/service"
	producgrpcerrror "github.com/MamangRust/simple_microservice_ecommerce/product/pkg/errors/grpc"
	"github.com/MamangRust/simple_microservice_ecommerce/product/pkg/logger"
	pbproduct "github.com/MamangRust/simple_microservice_ecommerce_pb/product"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type productCommandHandleGrpc struct {
	pbproduct.UnimplementedProductCommandServiceServer
	productCommandService service.ProductCommandService
	logger                logger.LoggerInterface
	mapper                productprotomapper.ProductCommandProtoMapper
}

func NewProductCommandHandleGrpc(Command service.ProductCommandService, logger logger.LoggerInterface) *productCommandHandleGrpc {
	return &productCommandHandleGrpc{
		productCommandService: Command,
		logger:                logger,
		mapper:                productprotomapper.NewProductCommandProtoMapper(),
	}
}

func (s *productCommandHandleGrpc) Create(ctx context.Context, request *pbproduct.CreateProductRequest) (*pbproduct.ApiResponseProduct, error) {
	req := &requests.CreateProductRequest{
		Name:  request.GetName(),
		Price: int(request.GetPrice()),
		Stock: int(request.GetStock()),
	}

	if err := req.Validate(); err != nil {
		s.logger.Warn("Invalid request data for Create", zap.Error(err), zap.String("product_name", request.GetName()))
		return nil, producgrpcerrror.ErrGrpcValidateCreateProduct
	}

	product, err := s.productCommandService.CreateProduct(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create product in service",
			zap.String("error_message", err.Message),
			zap.String("product_name", request.GetName()),
		)
		return nil, err.ToGRPCError()
	}

	s.logger.Info("Product successfully created",
		zap.Int("product_id", product.ID),
		zap.String("product_name", product.Name),
	)

	so := s.mapper.ToProtoResponseProduct("success", "Successfully created product", product)
	return so, nil
}

func (s *productCommandHandleGrpc) Update(ctx context.Context, request *pbproduct.UpdateProductRequest) (*pbproduct.ApiResponseProduct, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Warn("Invalid ProductID received on update request", zap.Int("product_id", id))
		return nil, producgrpcerrror.ErrGrpcInvalidID
	}

	req := &requests.UpdateProductRequest{
		ProductID: &id,
		Name:      request.GetName(),
		Price:     int(request.GetPrice()),
		Stock:     int(request.GetStock()),
	}

	if err := req.Validate(); err != nil {
		s.logger.Warn("Invalid request data for Update", zap.Error(err), zap.Int("product_id", id))
		return nil, producgrpcerrror.ErrGrpcValidateUpdateProduct
	}

	product, err := s.productCommandService.UpdateProduct(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update product in service",
			zap.String("error_message", err.Message),
			zap.Int("product_id", id),
		)
		return nil, err.ToGRPCError()
	}

	s.logger.Info("Product successfully updated", zap.Int("product_id", id))

	so := s.mapper.ToProtoResponseProduct("success", "Successfully updated product", product)
	return so, nil
}

func (s *productCommandHandleGrpc) UpdateProductCountStock(ctx context.Context, request *pbproduct.UpdateProductStockRequest) (*pbproduct.ApiResponseProduct, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Warn("Invalid ProductID received on stock update request", zap.Int("product_id", id))
		return nil, producgrpcerrror.ErrGrpcInvalidID
	}

	req := &requests.UpdateProductStockRequest{
		ProductID: id,
		Stock:     int(request.GetStock()),
	}

	product, err := s.productCommandService.UpdateProductStock(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update product stock in service",
			zap.String("error_message", err.Message),
			zap.Int("product_id", id),
		)
		return nil, err.ToGRPCError()
	}

	s.logger.Info("Product stock successfully updated",
		zap.Int("product_id", id),
		zap.Int("updated_stock", product.Stock),
	)

	so := s.mapper.ToProtoResponseProduct("success", "Successfully updated product stock", product)
	return so, nil
}

func (s *productCommandHandleGrpc) Trashed(ctx context.Context, request *pbproduct.FindByIdProductRequest) (*pbproduct.ApiResponseProductDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Warn("Invalid ProductID received on trash request", zap.Int("product_id", id))
		return nil, producgrpcerrror.ErrGrpcInvalidID
	}

	product, err := s.productCommandService.TrashedProduct(ctx, id)
	if err != nil {
		s.logger.Error("Failed to trash product in service",
			zap.String("error_message", err.Message),
			zap.Int("product_id", id),
		)
		return nil, err.ToGRPCError()
	}

	s.logger.Info("Product successfully trashed", zap.Int("product_id", id))

	so := s.mapper.ToProtoResponseProductDeleteAt("success", "Successfully trashed product", product)
	return so, nil
}

func (s *productCommandHandleGrpc) Restore(ctx context.Context, request *pbproduct.FindByIdProductRequest) (*pbproduct.ApiResponseProductDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Warn("Invalid ProductID received on restore request", zap.Int("product_id", id))
		return nil, producgrpcerrror.ErrGrpcInvalidID
	}

	product, err := s.productCommandService.RestoreProduct(ctx, id)
	if err != nil {
		s.logger.Error("Failed to restore product in service",
			zap.String("error_message", err.Message),
			zap.Int("product_id", id),
		)
		return nil, err.ToGRPCError()
	}

	s.logger.Info("Product successfully restored", zap.Int("product_id", id))

	so := s.mapper.ToProtoResponseProductDeleteAt("success", "Successfully restored product", product)
	return so, nil
}

func (s *productCommandHandleGrpc) DeleteProductPermanent(ctx context.Context, request *pbproduct.FindByIdProductRequest) (*pbproduct.ApiResponseProductDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Warn("Invalid ProductID received on permanent delete request", zap.Int("product_id", id))
		return nil, producgrpcerrror.ErrGrpcInvalidID
	}

	_, err := s.productCommandService.DeleteProductPermanent(ctx, id)
	if err != nil {
		s.logger.Error("Failed to permanently delete product in service",
			zap.String("error_message", err.Message),
			zap.Int("product_id", id),
		)
		return nil, err.ToGRPCError()
	}

	s.logger.Info("Product successfully permanently deleted", zap.Int("product_id", id))

	so := s.mapper.ToProtoResponseProductDelete("success", "Successfully deleted product permanently")
	return so, nil
}

func (s *productCommandHandleGrpc) RestoreAllProduct(ctx context.Context, _ *emptypb.Empty) (*pbproduct.ApiResponseProductAll, error) {
	_, err := s.productCommandService.RestoreAllProduct(ctx)
	if err != nil {
		s.logger.Error("Failed to restore all products in service", zap.String("error_message", err.Message))
		return nil, err.ToGRPCError()
	}

	s.logger.Info("All trashed products successfully restored")

	so := s.mapper.ToProtoResponseProductAll("success", "Successfully restored all products")
	return so, nil
}

func (s *productCommandHandleGrpc) DeleteAllProductPermanent(ctx context.Context, _ *emptypb.Empty) (*pbproduct.ApiResponseProductAll, error) {
	_, err := s.productCommandService.DeleteAllProductPermanent(ctx)
	if err != nil {

		s.logger.Error("Failed to permanently delete all products in service", zap.String("error_message", err.Message))
		return nil, err.ToGRPCError()
	}

	s.logger.Info("All trashed products successfully permanently deleted")

	so := s.mapper.ToProtoResponseProductAll("success", "Successfully deleted all products permanently")
	return so, nil
}
