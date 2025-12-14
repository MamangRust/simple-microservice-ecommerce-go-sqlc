package orderitemhandler

import (
	"context"
	"math"

	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/requests"
	orderitemprotomapper "github.com/MamangRust/simple_microservice_ecommerce/order/internal/mapper/proto/orderitem"
	orderitemservice "github.com/MamangRust/simple_microservice_ecommerce/order/internal/service/orderitem"
	orderitemgrpcerror "github.com/MamangRust/simple_microservice_ecommerce/order/pkg/errors/orderitem_errors/grpc"
	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/logger"
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb"
	pborderitem "github.com/MamangRust/simple_microservice_ecommerce_pb/order_item"
	"go.uber.org/zap"
)

type orderItemHandleGrpc struct {
	pborderitem.UnimplementedOrderItemServiceServer
	orderItemService orderitemservice.OrderItemQueryService
	logger           logger.LoggerInterface
	mapper           orderitemprotomapper.OrderItemProtoMapper
}

func NewOrderItemHandleGrpc(orderItem orderitemservice.OrderItemQueryService, logger logger.LoggerInterface) OrderItemHandleGrpc {
	return &orderItemHandleGrpc{
		orderItemService: orderItem,
		mapper:           orderitemprotomapper.NewOrderItemProtoMapper(),
		logger:           logger,
	}
}

func (s *orderItemHandleGrpc) FindAll(ctx context.Context, request *pborderitem.FindAllOrderItemRequest) (*pborderitem.ApiResponsePaginationOrderItem, error) {
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

	reqService := &requests.FindAllOrderItems{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orderItems, totalRecords, err := s.orderItemService.FindAllOrderItems(ctx, reqService)
	if err != nil {
		s.logger.Error("Failed to find all order items from service",
			zap.String("error_message", err.Message),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
		return nil, err.ToGRPCError()
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.Pagination{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("Successfully fetched order items",
		zap.Int("items_returned", len(orderItems)),
		zap.Int64("total_records", int64(*totalRecords)),
	)

	so := s.mapper.ToProtoResponsePaginationOrderItem(paginationMeta, "success", "Successfully fetched order items", orderItems)
	return so, nil
}

func (s *orderItemHandleGrpc) FindByActive(ctx context.Context, request *pborderitem.FindAllOrderItemRequest) (*pborderitem.ApiResponsePaginationOrderItemDeleteAt, error) {
	s.logger.Info("FindByActive handler called",
		zap.String("handler", "FindByActive"),
		zap.Int32("page", request.GetPage()),
		zap.Int32("pageSize", request.GetPageSize()),
		zap.String("search", request.GetSearch()),
	)

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

	reqService := &requests.FindAllOrderItems{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orderItems, totalRecords, err := s.orderItemService.FindByActive(ctx, reqService)
	if err != nil {
		s.logger.Error("Failed to find active order items from service",
			zap.String("error_message", err.Message),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
		return nil, err.ToGRPCError()
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.Pagination{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("Successfully fetched active order items",
		zap.Int("items_returned", len(orderItems)),
		zap.Int64("total_records", int64(*totalRecords)),
	)

	so := s.mapper.ToProtoResponsePaginationOrderItemDeleteAt(paginationMeta, "success", "Successfully fetched active order items", orderItems)
	return so, nil
}

func (s *orderItemHandleGrpc) FindByTrashed(ctx context.Context, request *pborderitem.FindAllOrderItemRequest) (*pborderitem.ApiResponsePaginationOrderItemDeleteAt, error) {
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

	reqService := &requests.FindAllOrderItems{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orderItems, totalRecords, err := s.orderItemService.FindByTrashed(ctx, reqService)
	if err != nil {
		s.logger.Error("Failed to find trashed order items from service",
			zap.String("error_message", err.Message),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
		return nil, err.ToGRPCError()
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.Pagination{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("Successfully fetched trashed order items",
		zap.Int("items_returned", len(orderItems)),
		zap.Int64("total_records", int64(*totalRecords)),
	)

	so := s.mapper.ToProtoResponsePaginationOrderItemDeleteAt(paginationMeta, "success", "Successfully fetched trashed order items", orderItems)
	return so, nil
}

func (s *orderItemHandleGrpc) FindOrderItemByOrder(ctx context.Context, request *pborderitem.FindByIdOrderItemRequest) (*pborderitem.ApiResponsesOrderItem, error) {
	orderId := int(request.GetId())

	if orderId == 0 {
		s.logger.Warn("Invalid OrderID received", zap.Int("orderId", orderId))
		return nil, orderitemgrpcerror.ErrGrpcInvalidID
	}

	orderItems, err := s.orderItemService.FindOrderItemByOrder(ctx, orderId)
	if err != nil {
		s.logger.Error("Failed to find order items by order ID from service",
			zap.String("error_message", err.Message),
			zap.Int("orderId", orderId),
		)
		return nil, err.ToGRPCError()
	}

	s.logger.Info("Successfully fetched order items for order",
		zap.Int("orderId", orderId),
		zap.Int("items_returned", len(orderItems)),
	)

	so := s.mapper.ToProtoResponsesOrderItem("success", "Successfully fetched order items for this order", orderItems)
	return so, nil
}
