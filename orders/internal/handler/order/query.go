package orderhandler

import (
	"context"
	"math"

	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/requests"
	orderprotomapper "github.com/MamangRust/simple_microservice_ecommerce/order/internal/mapper/proto/order"
	orderservice "github.com/MamangRust/simple_microservice_ecommerce/order/internal/service/order"
	ordergrpcerrors "github.com/MamangRust/simple_microservice_ecommerce/order/pkg/errors/order_errors/grpc"
	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/logger"
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb"
	pborder "github.com/MamangRust/simple_microservice_ecommerce_pb/order"
)

type orderQueryHandleGrpc struct {
	pborder.UnimplementedOrderQueryServiceServer
	orderQueryService orderservice.OrderQueryService
	logger            logger.LoggerInterface
	mapper            orderprotomapper.OrderQueryProtoMapper
}

func NewOrderQueryHandleGrpc(query orderservice.OrderQueryService, logger logger.LoggerInterface) *orderQueryHandleGrpc {
	return &orderQueryHandleGrpc{
		orderQueryService: query,
		logger:            logger,
		mapper:            orderprotomapper.NewOrderQueryProtoMapper(),
	}
}

func (s *orderQueryHandleGrpc) FindAll(ctx context.Context, request *pborder.FindAllOrderRequest) (*pborder.ApiResponsePaginationOrder, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := &requests.FindAllOrder{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orders, totalRecords, err := s.orderQueryService.FindAll(ctx, reqService)
	if err != nil {
		return nil, err.ToGRPCError()
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.Pagination{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapper.ToProtoResponsePaginationOrder(paginationMeta, "success", "Successfully fetched orders", orders)
	return so, nil
}

func (s *orderQueryHandleGrpc) FindById(ctx context.Context, request *pborder.FindByIdOrderRequest) (*pborder.ApiResponseOrder, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, ordergrpcerrors.ErrGrpcFailedInvalidId
	}

	order, err := s.orderQueryService.FindById(ctx, id)
	if err != nil {
		return nil, err.ToGRPCError()
	}

	so := s.mapper.ToProtoResponseOrder("success", "Successfully fetched order", order)
	return so, nil
}

func (s *orderQueryHandleGrpc) FindByActive(ctx context.Context, request *pborder.FindAllOrderRequest) (*pborder.ApiResponsePaginationOrderDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := &requests.FindAllOrder{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orders, totalRecords, err := s.orderQueryService.FindByActive(ctx, reqService)
	if err != nil {
		return nil, err.ToGRPCError()
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.Pagination{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapper.ToProtoResponsePaginationOrderDeleteAt(paginationMeta, "success", "Successfully fetched active orders", orders)
	return so, nil
}

func (s *orderQueryHandleGrpc) FindByTrashed(ctx context.Context, request *pborder.FindAllOrderRequest) (*pborder.ApiResponsePaginationOrderDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := &requests.FindAllOrder{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orders, totalRecords, err := s.orderQueryService.FindByTrashed(ctx, reqService)
	if err != nil {
		return nil, err.ToGRPCError()
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.Pagination{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapper.ToProtoResponsePaginationOrderDeleteAt(paginationMeta, "success", "Successfully fetched trashed orders", orders)
	return so, nil
}
