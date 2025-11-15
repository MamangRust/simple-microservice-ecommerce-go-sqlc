package orderservice

import (
	"context"
	"fmt"

	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
	grpcclient "github.com/MamangRust/simple_microservice_ecommerce/order/internal/grpc_client"
	orderservicemapper "github.com/MamangRust/simple_microservice_ecommerce/order/internal/mapper/response/service/order"
	orderrepository "github.com/MamangRust/simple_microservice_ecommerce/order/internal/repository/order"
	orderitemrepository "github.com/MamangRust/simple_microservice_ecommerce/order/internal/repository/order_item"
	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
}

func NewOrderCommandService(deps *OrderCommandDeps) OrderCommandService {
	return &orderCommandService{
		orderQueryRepository:       deps.orderQueryRepository,
		orderCommandRepository:     deps.orderCommandRepository,
		orderItemQueryRepository:   deps.orderItemQueryRepository,
		orderItemCommandRepository: deps.orderItemCommandRepository,
		mapper:                     deps.mapper,
		userGrpcClient:             deps.userGrpcClient,
		productGrpcClient:          deps.productGrpcClient,
		logger:                     deps.logger,
	}
}

func (s *orderCommandService) CreateOrder(ctx context.Context, req *requests.CreateOrderRequest) (*response.OrderResponse, *response.ErrorResponse) {
	s.logger.Info("Creating new order", zap.Int("user_id", req.UserID), zap.Int("items_count", len(req.Items)))

	_, errResp := s.userGrpcClient.FindById(ctx, int32(req.UserID))
	if errResp != nil {
		s.logger.Error("User not found when creating order", zap.Int("user_id", req.UserID), zap.String("error", errResp.Message))
		st, ok := status.FromError(errResp.ToGRPCError())
		if ok && st.Code() == codes.NotFound {
			return nil, response.NewErrorResponse("user not found", 404)
		}
		return nil, response.NewErrorResponse("failed to communicate with user service", 503)
	}

	order, err := s.orderCommandRepository.CreateOrder(ctx, &requests.CreateOrderRecordRequest{
		UserID:     req.UserID,
		TotalPrice: 0,
	})
	if err != nil {
		s.logger.Error("Failed to create order", zap.Int("user_id", req.UserID), zap.Error(err))
		return nil, response.NewErrorResponse("failed to create order", 500)
	}

	s.logger.Info("Order created successfully", zap.Int("order_id", order.ID))

	for _, item := range req.Items {
		product, errResp := s.productGrpcClient.FindById(ctx, int32(item.ProductID))
		if errResp != nil {
			s.logger.Error("Product not found when creating order", zap.Int("product_id", item.ProductID), zap.Int("order_id", order.ID), zap.String("error", errResp.Message))
			st, ok := status.FromError(errResp.ToGRPCError())
			if ok && st.Code() == codes.NotFound {
				return nil, response.NewErrorResponse(fmt.Sprintf("product with ID %d not found", item.ProductID), 400)
			}
			return nil, response.NewErrorResponse("failed to communicate with product service", 503)
		}

		if product.Data.Stock < item.Quantity {
			s.logger.Warn("Insufficient stock for product", zap.String("product_name", product.Data.Name), zap.Int("requested", item.Quantity), zap.Int("available", product.Data.Stock))
			return nil, response.NewErrorResponse(fmt.Sprintf("insufficient stock for product '%s'", product.Data.Name), 400)
		}

		_, err = s.orderItemCommandRepository.CreateOrderItem(ctx, &requests.CreateOrderItemRecordRequest{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Data.Price,
		})
		if err != nil {
			s.logger.Error("Failed to create order item", zap.Int("order_id", order.ID), zap.Int("product_id", item.ProductID), zap.Error(err))
			return nil, response.NewErrorResponse("failed to add item to order", 500)
		}

		product.Data.Stock -= item.Quantity
		_, errResp = s.productGrpcClient.UpdateProductStock(ctx, &requests.UpdateProductStockRequest{
			ProductID: product.Data.ID,
			Stock:     product.Data.Stock,
		})

		if errResp != nil {
			s.logger.Error("Failed to update product stock", zap.Int("product_id", product.Data.ID), zap.Int("new_stock", product.Data.Stock), zap.String("error", errResp.Message))
			st, ok := status.FromError(errResp.ToGRPCError())
			if ok && st.Code() == codes.NotFound {
				return nil, response.NewErrorResponse(fmt.Sprintf("product with ID %d not found during stock update", product.Data.ID), 404)
			}
			return nil, response.NewErrorResponse("failed to update product stock", 503)
		}

		s.logger.Info("Order item created and stock updated", zap.Int("order_id", order.ID), zap.Int("product_id", item.ProductID), zap.Int("quantity", item.Quantity), zap.Float64("price", float64(product.Data.Price)))
	}

	totalPrice, err := s.orderItemQueryRepository.CalculateTotalPrice(ctx, order.ID)
	if err != nil {
		s.logger.Error("Failed to calculate total price", zap.Int("order_id", order.ID), zap.Error(err))
		return nil, response.NewErrorResponse("failed to calculate total price", 500)
	}

	_, err = s.orderCommandRepository.UpdateOrder(ctx, &requests.UpdateOrderRecordRequest{
		OrderID:    order.ID,
		UserID:     req.UserID,
		TotalPrice: int(*totalPrice),
	})
	if err != nil {
		s.logger.Error("Failed to finalize order", zap.Int("order_id", order.ID), zap.Int("total_price", int(*totalPrice)), zap.Error(err))
		return nil, response.NewErrorResponse("failed to finalize order", 500)
	}

	s.logger.Info("Order created successfully", zap.Int("order_id", order.ID), zap.Float64("total_price", float64(*totalPrice)))

	so := s.mapper.ToOrderResponse(order)
	return so, nil
}

func (s *orderCommandService) UpdateOrder(ctx context.Context, req *requests.UpdateOrderRequest) (*response.OrderResponse, *response.ErrorResponse) {
	s.logger.Info("Updating order", zap.Int("order_id", req.OrderID), zap.Int("user_id", req.UserID), zap.Int("items_count", len(req.Items)))

	existingOrder, err := s.orderQueryRepository.FindById(ctx, req.OrderID)
	if err != nil {
		s.logger.Error("Order not found for update", zap.Int("order_id", req.OrderID), zap.Error(err))
		return nil, response.NewErrorResponse("order not found", 404)
	}

	_, errResp := s.userGrpcClient.FindById(ctx, int32(req.UserID))
	if errResp != nil {
		s.logger.Error("User not found when updating order", zap.Int("user_id", req.UserID), zap.String("error", errResp.Message))
		st, ok := status.FromError(errResp.ToGRPCError())
		if ok && st.Code() == codes.NotFound {
			return nil, response.NewErrorResponse("user not found", 404)
		}
		return nil, response.NewErrorResponse("failed to communicate with user service", 503)
	}

	for _, item := range req.Items {
		product, errResp := s.productGrpcClient.FindById(ctx, int32(item.ProductID))
		if errResp != nil {
			s.logger.Error("Product not found when updating order", zap.Int("product_id", item.ProductID), zap.Int("order_id", req.OrderID), zap.String("error", errResp.Message))
			st, ok := status.FromError(errResp.ToGRPCError())
			if ok && st.Code() == codes.NotFound {
				return nil, response.NewErrorResponse(fmt.Sprintf("product with ID %d not found", item.ProductID), 400)
			}
			return nil, response.NewErrorResponse("failed to communicate with product service", 503)
		}

		if item.OrderItemID > 0 {
			_, err := s.orderItemCommandRepository.UpdateOrderItem(ctx, &requests.UpdateOrderItemRecordRequest{
				OrderItemID: item.OrderItemID,
				ProductID:   item.ProductID,
				Quantity:    item.Quantity,
				Price:       product.Data.Price,
			})
			if err != nil {
				s.logger.Error("Failed to update order item", zap.Int("order_item_id", item.OrderItemID), zap.Int("product_id", item.ProductID), zap.Error(err))
				return nil, response.NewErrorResponse("failed to update order item", 500)
			}
			s.logger.Info("Order item updated", zap.Int("order_item_id", item.OrderItemID), zap.Int("product_id", item.ProductID), zap.Int("quantity", item.Quantity))
		} else {
			if product.Data.Stock < item.Quantity {
				s.logger.Warn("Insufficient stock for new order item", zap.String("product_name", product.Data.Name), zap.Int("requested", item.Quantity), zap.Int("available", product.Data.Stock))
				return nil, response.NewErrorResponse(fmt.Sprintf("insufficient stock for product '%s'", product.Data.Name), 400)
			}

			_, err = s.orderItemCommandRepository.CreateOrderItem(ctx, &requests.CreateOrderItemRecordRequest{
				OrderID:   req.OrderID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     product.Data.Price,
			})
			if err != nil {
				s.logger.Error("Failed to add new item to order", zap.Int("order_id", req.OrderID), zap.Int("product_id", item.ProductID), zap.Error(err))
				return nil, response.NewErrorResponse("failed to add new item to order", 500)
			}

			product.Data.Stock -= item.Quantity
			_, errResp = s.productGrpcClient.UpdateProductStock(ctx, &requests.UpdateProductStockRequest{
				ProductID: product.Data.ID,
				Stock:     product.Data.Stock,
			})

			if errResp != nil {
				s.logger.Error("Failed to update product stock for new item", zap.Int("product_id", product.Data.ID), zap.Int("new_stock", product.Data.Stock), zap.String("error", errResp.Message))
				st, ok := status.FromError(errResp.ToGRPCError())
				if ok && st.Code() == codes.NotFound {
					return nil, response.NewErrorResponse(fmt.Sprintf("product with ID %d not found during stock update", product.Data.ID), 404)
				}
				return nil, response.NewErrorResponse("failed to update product stock", 503)
			}

			s.logger.Info("New order item created and stock updated", zap.Int("order_id", req.OrderID), zap.Int("product_id", item.ProductID), zap.Int("quantity", item.Quantity), zap.Int("price", product.Data.Price))
		}
	}

	totalPrice, err := s.orderItemQueryRepository.CalculateTotalPrice(ctx, existingOrder.ID)
	if err != nil {
		s.logger.Error("Failed to calculate total price for updated order", zap.Int("order_id", existingOrder.ID), zap.Error(err))
		return nil, response.NewErrorResponse("failed to calculate total price", 500)
	}

	res, err := s.orderCommandRepository.UpdateOrder(ctx, &requests.UpdateOrderRecordRequest{
		OrderID:    existingOrder.ID,
		UserID:     req.UserID,
		TotalPrice: int(*totalPrice),
	})
	if err != nil {
		s.logger.Error("Failed to finalize order update", zap.Int("order_id", existingOrder.ID), zap.Int("total_price", int(*totalPrice)), zap.Error(err))
		return nil, response.NewErrorResponse("failed to finalize order update", 500)
	}

	s.logger.Info("Order updated successfully", zap.Int("order_id", existingOrder.ID), zap.Int("total_price", int(*totalPrice)))

	so := s.mapper.ToOrderResponse(res)
	return so, nil
}

func (s *orderCommandService) TrashedOrder(ctx context.Context, Order_id int) (*response.OrderResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Info("Trashing order", zap.Int("order_id", Order_id))

	res, err := s.orderCommandRepository.TrashedOrder(ctx, Order_id)
	if err != nil {
		s.logger.Error("Failed to trash order", zap.Int("order_id", Order_id), zap.Error(err))
		return nil, response.NewErrorResponse("failed to trash Order", 500)
	}

	s.logger.Info("Order trashed successfully", zap.Int("order_id", Order_id))

	so := s.mapper.ToOrderResponseDeleteAt(res)
	return so, nil
}

func (s *orderCommandService) RestoreOrder(ctx context.Context, Order_id int) (*response.OrderResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Info("Restoring order", zap.Int("order_id", Order_id))

	res, err := s.orderCommandRepository.RestoreOrder(ctx, Order_id)
	if err != nil {
		s.logger.Error("Failed to restore order", zap.Int("order_id", Order_id), zap.Error(err))
		return nil, response.NewErrorResponse("failed to restore Order", 500)
	}

	s.logger.Info("Order restored successfully", zap.Int("order_id", Order_id))

	so := s.mapper.ToOrderResponseDeleteAt(res)
	return so, nil
}

func (s *orderCommandService) DeleteOrderPermanent(ctx context.Context, Order_id int) (bool, *response.ErrorResponse) {
	s.logger.Info("Permanently deleting order", zap.Int("order_id", Order_id))

	_, err := s.orderCommandRepository.DeleteOrderPermanent(ctx, Order_id)
	if err != nil {
		s.logger.Error("Failed to delete order permanently", zap.Int("order_id", Order_id), zap.Error(err))
		return false, response.NewErrorResponse("failed to delete Order permanently", 500)
	}

	s.logger.Info("Order deleted permanently", zap.Int("order_id", Order_id))

	return true, nil
}

func (s *orderCommandService) RestoreAllOrder(ctx context.Context) (bool, *response.ErrorResponse) {
	s.logger.Info("Restoring all orders")

	_, err := s.orderCommandRepository.RestoreAllOrder(ctx)
	if err != nil {
		s.logger.Error("Failed to restore all orders", zap.Error(err))
		return false, response.NewErrorResponse("failed to restore all Orders", 500)
	}

	s.logger.Info("All orders restored successfully")

	return true, nil
}

func (s *orderCommandService) DeleteAllOrderPermanent(ctx context.Context) (bool, *response.ErrorResponse) {
	s.logger.Info("Permanently deleting all orders")

	_, err := s.orderCommandRepository.DeleteAllOrderPermanent(ctx)
	if err != nil {
		s.logger.Error("Failed to delete all orders permanently", zap.Error(err))
		return false, response.NewErrorResponse("failed to delete all Orders permanently", 500)
	}

	s.logger.Info("All orders deleted permanently")

	return true, nil
}
