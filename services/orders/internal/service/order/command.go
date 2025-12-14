package orderservice

import (
	"context"
	"fmt"

	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/errorhandler"
	grpcclient "github.com/MamangRust/simple_microservice_ecommerce/order/internal/grpc_client"
	orderservicemapper "github.com/MamangRust/simple_microservice_ecommerce/order/internal/mapper/response/service/order"
	mencache "github.com/MamangRust/simple_microservice_ecommerce/order/internal/redis"
	orderrepository "github.com/MamangRust/simple_microservice_ecommerce/order/internal/repository/order"
	orderitemrepository "github.com/MamangRust/simple_microservice_ecommerce/order/internal/repository/order_item"
	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/logger"
	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type OrderCommandDeps struct {
	orderQueryRepository       orderrepository.OrderQueryRepository
	orderCommandRepository     orderrepository.OrderCommandRepository
	orderItemQueryRepository   orderitemrepository.OrderItemQueryRepository
	orderItemCommandRepository orderitemrepository.OrderItemCommandRepository
	mapper                     orderservicemapper.OrderCommandResponseMapper
	userGrpcClient             grpcclient.UserGrpcClientHandler
	productGrpcClient          grpcclient.ProductGrpcClientHandler
	logger                     logger.LoggerInterface
	errorhandler               errorhandler.OrderCommandError
	mencache                   mencache.OrderCommandCache
}

type orderCommandService struct {
	orderQueryRepository       orderrepository.OrderQueryRepository
	orderCommandRepository     orderrepository.OrderCommandRepository
	orderItemQueryRepository   orderitemrepository.OrderItemQueryRepository
	orderItemCommandRepository orderitemrepository.OrderItemCommandRepository
	mapper                     orderservicemapper.OrderCommandResponseMapper
	userGrpcClient             grpcclient.UserGrpcClientHandler
	productGrpcClient          grpcclient.ProductGrpcClientHandler
	logger                     logger.LoggerInterface
	errorhandler               errorhandler.OrderCommandError
	mencache                   mencache.OrderCommandCache
	observability              observability.TraceLoggerObservability
}

func NewOrderCommandService(deps *OrderCommandDeps) OrderCommandService {
	observability, _ := observability.NewObservability("order-command-service", deps.logger)

	return &orderCommandService{
		orderQueryRepository:       deps.orderQueryRepository,
		orderCommandRepository:     deps.orderCommandRepository,
		orderItemQueryRepository:   deps.orderItemQueryRepository,
		orderItemCommandRepository: deps.orderItemCommandRepository,
		mapper:                     deps.mapper,
		userGrpcClient:             deps.userGrpcClient,
		productGrpcClient:          deps.productGrpcClient,
		logger:                     deps.logger,
		errorhandler:               deps.errorhandler,
		mencache:                   deps.mencache,
		observability:              observability,
	}
}

func (s *orderCommandService) CreateOrder(ctx context.Context, req *requests.CreateOrderRequest) (*response.OrderResponse, *response.ErrorResponse) {
	const method = "CreateOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("user_id", req.UserID), attribute.Int("items_count", len(req.Items)))
	defer func() {
		end(status)
	}()

	_, errResp := s.userGrpcClient.FindById(ctx, int32(req.UserID))
	if errResp != nil {
		status = "error"
		return s.errorhandler.HandleCreateOrderError(fmt.Errorf("user not found: %s", errResp.Message), method, "FAILED_USER_NOT_FOUND", span, &status, zap.Int("user_id", req.UserID))
	}

	order, err := s.orderCommandRepository.CreateOrder(ctx, &requests.CreateOrderRecordRequest{
		UserID:     req.UserID,
		TotalPrice: 0,
	})
	if err != nil {
		status = "error"
		return s.errorhandler.HandleCreateOrderError(err, method, "FAILED_CREATE_ORDER", span, &status, zap.Int("user_id", req.UserID))
	}

	for _, item := range req.Items {
		product, errResp := s.productGrpcClient.FindById(ctx, int32(item.ProductID))
		if errResp != nil {
			status = "error"
			return s.errorhandler.HandleCreateOrderError(fmt.Errorf("product not found: %s", errResp.Message), method, "FAILED_PRODUCT_NOT_FOUND", span, &status, zap.Int("product_id", item.ProductID), zap.Int("order_id", order.ID))
		}

		if product.Data.Stock < item.Quantity {
			status = "error"
			errMsg := fmt.Sprintf("insufficient stock for product '%s'", product.Data.Name)
			return s.errorhandler.HandleErrorInsufficientStockTemplate(fmt.Errorf("%s", errMsg), method, "FAILED_INSUFFICIENT_STOCK", span, &status, response.NewErrorResponse(errMsg, 400), zap.String("product_name", product.Data.Name), zap.Int("requested", item.Quantity), zap.Int("available", product.Data.Stock))
		}

		_, err = s.orderItemCommandRepository.CreateOrderItem(ctx, &requests.CreateOrderItemRecordRequest{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Data.Price,
		})
		if err != nil {
			status = "error"
			return s.errorhandler.HandleCreateOrderError(err, method, "FAILED_CREATE_ORDER_ITEM", span, &status, zap.Int("order_id", order.ID), zap.Int("product_id", item.ProductID))
		}

		product.Data.Stock -= item.Quantity
		_, errResp = s.productGrpcClient.UpdateProductStock(ctx, &requests.UpdateProductStockRequest{
			ProductID: product.Data.ID,
			Stock:     product.Data.Stock,
		})

		if errResp != nil {
			status = "error"
			return s.errorhandler.HandleCreateOrderError(fmt.Errorf("failed to update product stock: %s", errResp.Message), method, "FAILED_UPDATE_PRODUCT_STOCK", span, &status, zap.Int("product_id", product.Data.ID))
		}
	}

	totalPrice, err := s.orderItemQueryRepository.CalculateTotalPrice(ctx, order.ID)
	if err != nil {
		status = "error"
		return s.errorhandler.HandleCreateOrderError(err, method, "FAILED_CALCULATE_TOTAL_PRICE", span, &status, zap.Int("order_id", order.ID))
	}

	_, err = s.orderCommandRepository.UpdateOrder(ctx, &requests.UpdateOrderRecordRequest{
		OrderID:    order.ID,
		UserID:     req.UserID,
		TotalPrice: int(*totalPrice),
	})
	if err != nil {
		status = "error"
		return s.errorhandler.HandleCreateOrderError(err, method, "FAILED_FINALIZE_ORDER", span, &status, zap.Int("order_id", order.ID))
	}

	so := s.mapper.ToOrderResponse(order)

	s.mencache.InvalidateAllOrders(ctx)
	s.mencache.InvalidateActiveOrders(ctx)

	logSuccess("Order created successfully", zap.Int("order_id", order.ID), zap.Float64("total_price", float64(*totalPrice)))

	return so, nil
}

func (s *orderCommandService) UpdateOrder(ctx context.Context, req *requests.UpdateOrderRequest) (*response.OrderResponse, *response.ErrorResponse) {
	const method = "UpdateOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("order_id", req.OrderID), attribute.Int("user_id", req.UserID))
	defer func() {
		end(status)
	}()

	existingOrder, err := s.orderQueryRepository.FindById(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return s.errorhandler.HandleUpdateOrderError(err, method, "FAILED_ORDER_NOT_FOUND", span, &status, zap.Int("order_id", req.OrderID))
	}

	_, errResp := s.userGrpcClient.FindById(ctx, int32(req.UserID))
	if errResp != nil {
		status = "error"
		return s.errorhandler.HandleUpdateOrderError(fmt.Errorf("user not found: %s", errResp.Message), method, "FAILED_USER_NOT_FOUND", span, &status, zap.Int("user_id", req.UserID))
	}

	for _, item := range req.Items {
		product, errResp := s.productGrpcClient.FindById(ctx, int32(item.ProductID))
		if errResp != nil {
			status = "error"
			return s.errorhandler.HandleUpdateOrderError(fmt.Errorf("product not found: %s", errResp.Message), method, "FAILED_PRODUCT_NOT_FOUND", span, &status, zap.Int("product_id", item.ProductID))
		}

		if item.OrderItemID > 0 {
			_, err := s.orderItemCommandRepository.UpdateOrderItem(ctx, &requests.UpdateOrderItemRecordRequest{
				OrderItemID: item.OrderItemID,
				ProductID:   item.ProductID,
				Quantity:    item.Quantity,
				Price:       product.Data.Price,
			})
			if err != nil {
				status = "error"
				return s.errorhandler.HandleUpdateOrderError(err, method, "FAILED_UPDATE_ORDER_ITEM", span, &status, zap.Int("order_item_id", item.OrderItemID))
			}
		} else {
			if product.Data.Stock < item.Quantity {
				status = "error"
				errMsg := fmt.Sprintf("insufficient stock for product '%s'", product.Data.Name)
				return s.errorhandler.HandleErrorInsufficientStockTemplate(fmt.Errorf("%s", errMsg), method, "FAILED_INSUFFICIENT_STOCK", span, &status, response.NewErrorResponse(errMsg, 400), zap.String("product_name", product.Data.Name), zap.Int("requested", item.Quantity), zap.Int("available", product.Data.Stock))
			}

			_, err = s.orderItemCommandRepository.CreateOrderItem(ctx, &requests.CreateOrderItemRecordRequest{
				OrderID:   req.OrderID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     product.Data.Price,
			})
			if err != nil {
				status = "error"
				return s.errorhandler.HandleUpdateOrderError(err, method, "_FAILED_CREATE_ORDER_ITEM", span, &status, zap.Int("order_id", req.OrderID), zap.Int("product_id", item.ProductID))
			}

			product.Data.Stock -= item.Quantity
			_, errResp = s.productGrpcClient.UpdateProductStock(ctx, &requests.UpdateProductStockRequest{
				ProductID: product.Data.ID,
				Stock:     product.Data.Stock,
			})

			if errResp != nil {
				status = "error"
				return s.errorhandler.HandleUpdateOrderError(fmt.Errorf("failed to update product stock: %s", errResp.Message), method, "FAILED_UPDATE_PRODUCT_STOCK", span, &status, zap.Int("product_id", product.Data.ID))
			}
		}
	}

	totalPrice, err := s.orderItemQueryRepository.CalculateTotalPrice(ctx, existingOrder.ID)
	if err != nil {
		status = "error"
		return s.errorhandler.HandleUpdateOrderError(err, method, "FAILED_CALCULATE_TOTAL_PRICE", span, &status, zap.Int("order_id", existingOrder.ID))
	}

	res, err := s.orderCommandRepository.UpdateOrder(ctx, &requests.UpdateOrderRecordRequest{
		OrderID:    existingOrder.ID,
		UserID:     req.UserID,
		TotalPrice: int(*totalPrice),
	})
	if err != nil {
		status = "error"
		return s.errorhandler.HandleUpdateOrderError(err, method, "FAILED_FINALIZE_ORDER_UPDATE", span, &status, zap.Int("order_id", existingOrder.ID))
	}

	so := s.mapper.ToOrderResponse(res)

	s.mencache.DeleteOrderCache(ctx, req.OrderID)
	s.mencache.InvalidateActiveOrders(ctx)
	s.mencache.InvalidateAllOrders(ctx)

	logSuccess("Order updated successfully", zap.Int("order_id", existingOrder.ID), zap.Int("total_price", int(*totalPrice)))

	return so, nil
}

func (s *orderCommandService) TrashedOrder(ctx context.Context, Order_id int) (*response.OrderResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashedOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("order_id", Order_id))
	defer func() {
		end(status)
	}()

	res, err := s.orderCommandRepository.TrashedOrder(ctx, Order_id)
	if err != nil {
		status = "error"
		return s.errorhandler.HandleTrashedOrderError(err, method, "FAILED_TRASH_ORDER", span, &status, zap.Int("order_id", Order_id))
	}

	so := s.mapper.ToOrderResponseDeleteAt(res)

	s.mencache.DeleteOrderCache(ctx, Order_id)
	s.mencache.InvalidateActiveOrders(ctx)
	s.mencache.InvalidateAllOrders(ctx)
	s.mencache.InvalidateTrashedOrders(ctx)

	logSuccess("Order trashed successfully", zap.Int("order_id", Order_id))

	return so, nil
}

func (s *orderCommandService) RestoreOrder(ctx context.Context, Order_id int) (*response.OrderResponseDeleteAt, *response.ErrorResponse) {
	const method = "RestoreOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("order_id", Order_id))
	defer func() {
		end(status)
	}()

	res, err := s.orderCommandRepository.RestoreOrder(ctx, Order_id)
	if err != nil {
		status = "error"
		return s.errorhandler.HandleRestoreOrderError(err, method, "FAILED_RESTORE_ORDER", span, &status, zap.Int("order_id", Order_id))
	}

	so := s.mapper.ToOrderResponseDeleteAt(res)

	s.mencache.DeleteOrderCache(ctx, Order_id)
	s.mencache.InvalidateActiveOrders(ctx)
	s.mencache.InvalidateAllOrders(ctx)
	s.mencache.InvalidateTrashedOrders(ctx)

	logSuccess("Order restored successfully", zap.Int("order_id", Order_id))

	return so, nil
}

func (s *orderCommandService) DeleteOrderPermanent(ctx context.Context, Order_id int) (bool, *response.ErrorResponse) {
	const method = "DeleteOrderPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("order_id", Order_id))
	defer func() {
		end(status)
	}()

	_, err := s.orderCommandRepository.DeleteOrderPermanent(ctx, Order_id)
	if err != nil {
		status = "error"
		return s.errorhandler.HandleDeleteOrderError(err, method, "FAILED_DELETE_ORDER_PERMANENT", span, &status, zap.Int("order_id", Order_id))
	}

	s.mencache.DeleteOrderCache(ctx, Order_id)
	s.mencache.InvalidateActiveOrders(ctx)
	s.mencache.InvalidateAllOrders(ctx)
	s.mencache.InvalidateTrashedOrders(ctx)

	logSuccess("Order deleted permanently", zap.Int("order_id", Order_id))

	return true, nil
}

func (s *orderCommandService) RestoreAllOrder(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "RestoreAllOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	_, err := s.orderCommandRepository.RestoreAllOrder(ctx)
	if err != nil {
		status = "error"
		return s.errorhandler.HandleRestoreAllOrderError(err, method, "FAILED_RESTORE_ALL_ORDERS", span, &status)
	}

	s.mencache.InvalidateActiveOrders(ctx)
	s.mencache.InvalidateAllOrders(ctx)
	s.mencache.InvalidateTrashedOrders(ctx)

	logSuccess("All orders restored successfully")

	return true, nil
}

func (s *orderCommandService) DeleteAllOrderPermanent(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "DeleteAllOrderPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	_, err := s.orderCommandRepository.DeleteAllOrderPermanent(ctx)
	if err != nil {
		status = "error"
		return s.errorhandler.HandleDeleteAllOrderError(err, method, "FAILED_DELETE_ALL_ORDERS_PERMANENT", span, &status)
	}

	s.mencache.InvalidateActiveOrders(ctx)
	s.mencache.InvalidateAllOrders(ctx)
	s.mencache.InvalidateTrashedOrders(ctx)

	logSuccess("All orders deleted permanently")

	return true, nil
}
