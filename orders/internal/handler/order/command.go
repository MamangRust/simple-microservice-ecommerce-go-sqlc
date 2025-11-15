package orderhandler

import (
	"context"

	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/requests"
	orderprotomapper "github.com/MamangRust/simple_microservice_ecommerce/order/internal/mapper/proto/order"
	orderservice "github.com/MamangRust/simple_microservice_ecommerce/order/internal/service/order"
	ordergrpcerrors "github.com/MamangRust/simple_microservice_ecommerce/order/pkg/errors/order_errors/grpc"
	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/logger"
	pborder "github.com/MamangRust/simple_microservice_ecommerce_pb/order"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type orderCommandHandleGrpc struct {
	pborder.UnimplementedOrderCommandServiceServer
	orderCommandService orderservice.OrderCommandService
	logger              logger.LoggerInterface
	mapper              orderprotomapper.OrderCommandProtoMapper
}

func NewOrderCommandHandleGrpc(command orderservice.OrderCommandService, logger logger.LoggerInterface) *orderCommandHandleGrpc {
	return &orderCommandHandleGrpc{
		orderCommandService: command,
		logger:              logger,
		mapper:              orderprotomapper.NewOrderCommandProtoMapper(),
	}
}

func (s *orderCommandHandleGrpc) Create(ctx context.Context, request *pborder.CreateOrderRequest) (*pborder.ApiResponseOrder, error) {
	if request.GetUserId() == 0 {
		s.logger.Warn("Invalid ID received on create order request",
			zap.Int32("userId", request.GetUserId()),
		)
		return nil, ordergrpcerrors.ErrGrpcFailedInvalidId
	}

	var items []requests.CreateOrderItemRequest
	for _, item := range request.GetItems() {
		if item.GetProductId() == 0 || item.GetQuantity() == 0 {
			s.logger.Warn("Invalid item data received on create order request",
				zap.Int32("productId", item.GetProductId()),
				zap.Int32("quantity", item.GetQuantity()),
			)
			return nil, ordergrpcerrors.ErrGrpcFailedInvalidId
		}
		items = append(items, requests.CreateOrderItemRequest{
			ProductID: int(item.GetProductId()),
			Quantity:  int(item.GetQuantity()),
			Price:     int(item.GetPrice()),
		})
	}

	reqService := &requests.CreateOrderRequest{
		UserID: int(request.GetUserId()),
		Items:  items,
	}

	order, err := s.orderCommandService.CreateOrder(ctx, reqService)
	if err != nil {
		s.logger.Error("Failed to create order in service",
			zap.String("error_message", err.Message),
			zap.Int("userId", reqService.UserID),
		)
		return nil, err.ToGRPCError()
	}

	s.logger.Info("Order successfully created",
		zap.Int("orderId", order.ID),
		zap.Int("userId", reqService.UserID),
	)

	so := s.mapper.ToProtoResponseOrder("success", "Successfully created order", order)
	return so, nil
}

func (s *orderCommandHandleGrpc) Update(ctx context.Context, request *pborder.UpdateOrderRequest) (*pborder.ApiResponseOrder, error) {
	if request.GetOrderId() == 0 || request.GetUserId() == 0 {
		s.logger.Warn("Invalid ID received on update order request",
			zap.Int32("orderId", request.GetOrderId()),
			zap.Int32("userId", request.GetUserId()),
		)
		return nil, ordergrpcerrors.ErrGrpcFailedInvalidId
	}

	var items []requests.UpdateOrderItemRequest
	for _, item := range request.GetItems() {
		if item.GetOrderItemId() == 0 || item.GetProductId() == 0 || item.GetQuantity() == 0 {
			s.logger.Warn("Invalid item data received on update order request",
				zap.Int32("orderItemId", item.GetOrderItemId()),
				zap.Int32("productId", item.GetProductId()),
				zap.Int32("quantity", item.GetQuantity()),
			)
			return nil, ordergrpcerrors.ErrGrpcFailedInvalidId
		}
		items = append(items, requests.UpdateOrderItemRequest{
			OrderItemID: int(item.GetOrderItemId()),
			ProductID:   int(item.GetProductId()),
			Quantity:    int(item.GetQuantity()),
			Price:       int(item.GetPrice()),
		})
	}

	orderId := int(request.GetOrderId())

	reqService := &requests.UpdateOrderRequest{
		OrderID: orderId,
		UserID:  int(request.GetUserId()),
		Items:   items,
	}

	order, err := s.orderCommandService.UpdateOrder(ctx, reqService)
	if err != nil {
		s.logger.Error("Failed to update order in service",
			zap.String("error_message", err.Message),
			zap.Int("orderId", orderId),
		)
		return nil, err.ToGRPCError()
	}

	s.logger.Info("Order successfully updated",
		zap.Int("orderId", orderId),
	)

	so := s.mapper.ToProtoResponseOrder("success", "Successfully updated order", order)
	return so, nil
}

func (s *orderCommandHandleGrpc) Trashed(ctx context.Context, request *pborder.FindByIdOrderRequest) (*pborder.ApiResponseOrderDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Warn("Invalid OrderID received on trash request", zap.Int("orderId", id))
		return nil, ordergrpcerrors.ErrGrpcFailedInvalidId
	}

	order, err := s.orderCommandService.TrashedOrder(ctx, id)
	if err != nil {
		s.logger.Error("Failed to trash order in service",
			zap.String("error_message", err.Message),
			zap.Int("orderId", id),
		)
		return nil, err.ToGRPCError()
	}

	s.logger.Info("Order successfully trashed",
		zap.Int("orderId", id),
	)

	so := s.mapper.ToProtoResponseOrderDeleteAt("success", "Successfully trashed order", order)
	return so, nil
}

func (s *orderCommandHandleGrpc) Restore(ctx context.Context, request *pborder.FindByIdOrderRequest) (*pborder.ApiResponseOrderDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Warn("Invalid OrderID received on restore request", zap.Int("orderId", id))
		return nil, ordergrpcerrors.ErrGrpcFailedInvalidId
	}

	order, err := s.orderCommandService.RestoreOrder(ctx, id)
	if err != nil {
		s.logger.Error("Failed to restore order in service",
			zap.String("error_message", err.Message),
			zap.Int("orderId", id),
		)
		return nil, err.ToGRPCError()
	}

	s.logger.Info("Order successfully restored",
		zap.Int("orderId", id),
	)

	so := s.mapper.ToProtoResponseOrderDeleteAt("success", "Successfully restored order", order)
	return so, nil
}

func (s *orderCommandHandleGrpc) DeleteOrderPermanent(ctx context.Context, request *pborder.FindByIdOrderRequest) (*pborder.ApiResponseOrderDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Warn("Invalid OrderID received on permanent delete request", zap.Int("orderId", id))
		return nil, ordergrpcerrors.ErrGrpcFailedInvalidId
	}

	_, err := s.orderCommandService.DeleteOrderPermanent(ctx, id)
	if err != nil {
		s.logger.Error("Failed to permanently delete order in service",
			zap.String("error_message", err.Message),
			zap.Int("orderId", id),
		)
		return nil, err.ToGRPCError()
	}

	s.logger.Info("Order successfully permanently deleted",
		zap.Int("orderId", id),
	)

	so := s.mapper.ToProtoResponseOrderDelete("success", "Successfully deleted order permanently")
	return so, nil
}

func (s *orderCommandHandleGrpc) RestoreAllOrder(ctx context.Context, _ *emptypb.Empty) (*pborder.ApiResponseOrderAll, error) {
	_, err := s.orderCommandService.RestoreAllOrder(ctx)
	if err != nil {
		s.logger.Error("Failed to restore all orders in service", zap.String("error_message", err.Message))
		return nil, err.ToGRPCError()
	}

	s.logger.Info("All trashed orders successfully restored")

	so := s.mapper.ToProtoResponseOrderAll("success", "Successfully restored all orders")
	return so, nil
}

func (s *orderCommandHandleGrpc) DeleteAllOrder(ctx context.Context, _ *emptypb.Empty) (*pborder.ApiResponseOrderAll, error) {
	_, err := s.orderCommandService.DeleteAllOrderPermanent(ctx)
	if err != nil {
		s.logger.Error("Failed to permanently delete all orders in service", zap.String("error_message", err.Message))
		return nil, err.ToGRPCError()
	}

	s.logger.Info("All trashed orders successfully permanently deleted")

	so := s.mapper.ToProtoResponseOrderAll("success", "Successfully deleted all orders permanently")
	return so, nil
}
