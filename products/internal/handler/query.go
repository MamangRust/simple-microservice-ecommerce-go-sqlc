package handler

import (
	"context"
	"math"

	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/requests"
	productprotomapper "github.com/MamangRust/simple_microservice_ecommerce/product/internal/mapper/proto/product"
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/service"
	producgrpcerrror "github.com/MamangRust/simple_microservice_ecommerce/product/pkg/errors/grpc"
	"github.com/MamangRust/simple_microservice_ecommerce/product/pkg/logger"
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb"
	pbproduct "github.com/MamangRust/simple_microservice_ecommerce_pb/product"
	"go.uber.org/zap"
)

type productQueryHandleGrpc struct {
	pbproduct.UnimplementedProductQueryServiceServer
	productQueryService service.ProductQueryService
	logger              logger.LoggerInterface
	mapper              productprotomapper.ProductQueryProtoMapper
}

func NewProductQueryHandleGrpc(query service.ProductQueryService, logger logger.LoggerInterface) *productQueryHandleGrpc {
	return &productQueryHandleGrpc{
		productQueryService: query,
		logger:              logger,
		mapper:              productprotomapper.NewProductQueryProtoMapper(),
	}
}

func (s *productQueryHandleGrpc) FindAll(ctx context.Context, request *pbproduct.FindAllProductRequest) (*pbproduct.ApiResponsePaginationProduct, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		s.logger.Warn("Invalid page number received, defaulting to 1", zap.Int("page", page))
		page = 1
	}

	if pageSize <= 0 {
		s.logger.Warn("Invalid pageSize received, defaulting to 10", zap.Int("pageSize", pageSize))
		pageSize = 10
	}

	reqService := &requests.FindAllProduct{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	products, totalRecords, err := s.productQueryService.FindAll(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to find all products", zap.String("error_message", err.Message))
		return nil, err.ToGRPCError()
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.Pagination{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("Successfully fetched products",
		zap.Int("count", len(products)),
		zap.Int64("total_records", int64(*totalRecords)),
	)

	so := s.mapper.ToProtoResponsePaginationProduct(paginationMeta, "success", "Successfully fetched products", products)
	return so, nil
}

func (s *productQueryHandleGrpc) FindById(ctx context.Context, request *pbproduct.FindByIdProductRequest) (*pbproduct.ApiResponseProduct, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Warn("Invalid ProductID received", zap.Int("product_id", id))
		return nil, producgrpcerrror.ErrGrpcInvalidID
	}

	product, err := s.productQueryService.FindByID(ctx, id)

	if err != nil {
		s.logger.Error("Failed to find product by ID", zap.String("error_message", err.Message), zap.Int("product_id", id))
		return nil, err.ToGRPCError()
	}

	s.logger.Info("Successfully fetched product", zap.Int("product_id", id))

	so := s.mapper.ToProtoResponseProduct("success", "Successfully fetched product", product)
	return so, nil
}

func (s *productQueryHandleGrpc) FindByActive(ctx context.Context, request *pbproduct.FindAllProductRequest) (*pbproduct.ApiResponsePaginationProductDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		s.logger.Warn("Invalid page number received, defaulting to 1", zap.Int("page", page))
		page = 1
	}
	if pageSize <= 0 {
		s.logger.Warn("Invalid pageSize received, defaulting to 10", zap.Int("pageSize", pageSize))
		pageSize = 10
	}

	reqService := &requests.FindAllProduct{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	products, totalRecords, err := s.productQueryService.FindByActive(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to find active products", zap.String("error_message", err.Message))
		return nil, err.ToGRPCError()
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.Pagination{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("Successfully fetched active products",
		zap.Int("count", len(products)),
		zap.Int64("total_records", int64(*totalRecords)),
	)

	so := s.mapper.ToProtoResponsePaginationProductDeleteAt(paginationMeta, "success", "Successfully fetched active products", products)
	return so, nil
}

func (s *productQueryHandleGrpc) FindByTrashed(ctx context.Context, request *pbproduct.FindAllProductRequest) (*pbproduct.ApiResponsePaginationProductDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		s.logger.Warn("Invalid page number received, defaulting to 1", zap.Int("page", page))
		page = 1
	}
	if pageSize <= 0 {
		s.logger.Warn("Invalid pageSize received, defaulting to 10", zap.Int("pageSize", pageSize))
		pageSize = 10
	}

	reqService := &requests.FindAllProduct{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	products, totalRecords, err := s.productQueryService.FindByTrashed(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to find trashed products", zap.String("error_message", err.Message))
		return nil, err.ToGRPCError()
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.Pagination{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("Successfully fetched trashed products",
		zap.Int("count", len(products)),
		zap.Int64("total_records", int64(*totalRecords)),
	)

	so := s.mapper.ToProtoResponsePaginationProductDeleteAt(paginationMeta, "success", "Successfully fetched trashed products", products)
	return so, nil
}
