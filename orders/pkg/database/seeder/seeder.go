package seeder

import (
	"context"
	"fmt"

	db "github.com/MamangRust/simple_microservice_ecommerce/order/pkg/database/schema"
	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/logger"
	"go.uber.org/zap"
)

type OrderSeeder struct {
	db     *db.Queries
	logger logger.LoggerInterface
}

func NewOrderSeeder(db *db.Queries, logger logger.LoggerInterface) *OrderSeeder {
	return &OrderSeeder{
		db:     db,
		logger: logger,
	}
}

func (s *OrderSeeder) SeedAll(ctx context.Context) error {
	s.logger.Info("[SEEDER] Starting orders service seeding...")

	if err := s.SeedOrders(ctx); err != nil {
		s.logger.Error("[SEEDER] Failed to seed orders", zap.Error(err))
		return fmt.Errorf("failed to seed orders: %w", err)
	}

	if err := s.SeedOrderItems(ctx); err != nil {
		s.logger.Error("[SEEDER] Failed to seed order items", zap.Error(err))
		return fmt.Errorf("failed to seed order items: %w", err)
	}

	s.logger.Info("[SEEDER] Orders service seeding completed successfully")
	return nil
}

func (s *OrderSeeder) SeedOrders(ctx context.Context) error {
	s.logger.Info("[SEEDER] Seeding orders...")

	orders := []db.CreateOrderParams{
		{UserID: 1, TotalPrice: 150000},
		{UserID: 2, TotalPrice: 75000},
		{UserID: 3, TotalPrice: 225000},
		{UserID: 4, TotalPrice: 50000},
		{UserID: 5, TotalPrice: 300000},
	}

	createdCount := 0
	for _, order := range orders {
		orderExists := false
		existingOrders, err := s.db.GetOrders(ctx, db.GetOrdersParams{Limit: 1000, Offset: 0})
		if err == nil {
			for _, existingOrder := range existingOrders {
				if existingOrder.UserID == order.UserID {
					orderExists = true
					break
				}
			}
		}

		if orderExists {
			s.logger.Info("[SEEDER] Order for user already exists, skipping", zap.Int32("user_id", order.UserID))
			continue
		}

		newOrder, err := s.db.CreateOrder(ctx, db.CreateOrderParams{
			UserID:     order.UserID,
			TotalPrice: order.TotalPrice,
		})

		if err != nil {
			s.logger.Error("[SEEDER] Error inserting order", zap.Error(err), zap.Int32("user_id", order.UserID))
			return err
		}

		s.logger.Info("[SEEDER] Successfully created order",
			zap.Int64("order_id", int64(newOrder.OrderID)),
			zap.Int32("user_id", order.UserID),
			zap.Int64("total_price", int64(order.TotalPrice)),
		)
		createdCount++
	}

	s.logger.Info("[SEEDER] Order seeding summary", zap.Int("total_orders_in_list", len(orders)), zap.Int("orders_created", createdCount))
	return nil
}

func (s *OrderSeeder) SeedOrderItems(ctx context.Context) error {
	s.logger.Info("[SEEDER] Seeding order items...")

	orderItems := []db.CreateOrderItemParams{
		{OrderID: 1, ProductID: 1, Quantity: 2, Price: 50000},
		{OrderID: 1, ProductID: 2, Quantity: 1, Price: 25000},
		{OrderID: 1, ProductID: 3, Quantity: 3, Price: 15000},
		{OrderID: 2, ProductID: 1, Quantity: 1, Price: 50000},
		{OrderID: 2, ProductID: 4, Quantity: 1, Price: 25000},
		{OrderID: 3, ProductID: 2, Quantity: 3, Price: 25000},
		{OrderID: 3, ProductID: 3, Quantity: 5, Price: 15000},
		{OrderID: 3, ProductID: 5, Quantity: 2, Price: 30000},
		{OrderID: 4, ProductID: 1, Quantity: 1, Price: 50000},
		{OrderID: 5, ProductID: 4, Quantity: 4, Price: 25000},
		{OrderID: 5, ProductID: 5, Quantity: 5, Price: 30000},
		{OrderID: 5, ProductID: 6, Quantity: 2, Price: 50000},
	}

	createdCount := 0
	for _, item := range orderItems {
		// Logika pengecekan sudah ada, kita hanya perlu menambahkan log
		itemExists := false
		existingItems, err := s.db.GetOrderItems(ctx, db.GetOrderItemsParams{Limit: 1000, Offset: 0})
		if err == nil {
			for _, existingItem := range existingItems {
				if existingItem.OrderID == item.OrderID && existingItem.ProductID == item.ProductID {
					itemExists = true
					break
				}
			}
		}

		if itemExists {
			s.logger.Info("[SEEDER] Order item already exists, skipping",
				zap.Int32("order_id", item.OrderID),
				zap.Int32("product_id", item.ProductID),
			)
			continue
		}

		_, err = s.db.CreateOrderItem(ctx, db.CreateOrderItemParams{
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})

		if err != nil {
			s.logger.Error("[SEEDER] Error inserting order item",
				zap.Error(err),
				zap.Int32("order_id", item.OrderID),
				zap.Int32("product_id", item.ProductID),
			)
			return err
		}

		s.logger.Info("[SEEDER] Successfully created order item",
			zap.Int32("order_id", item.OrderID),
			zap.Int32("product_id", item.ProductID),
			zap.Int32("quantity", item.Quantity),
			zap.Int64("price", int64(item.Price)),
		)
		createdCount++
	}

	s.logger.Info("[SEEDER] Order item seeding summary", zap.Int("total_items_in_list", len(orderItems)), zap.Int("items_created", createdCount))
	return nil
}

func (s *OrderSeeder) ClearAll(ctx context.Context) error {
	s.logger.Info("[SEEDER] Clearing orders service data...")

	deletedItemCount := 0
	orderItems, err := s.db.GetOrderItems(ctx, db.GetOrderItemsParams{Limit: 1000, Offset: 0})
	if err == nil {
		for _, item := range orderItems {
			err = s.db.DeleteOrderItemPermanently(ctx, item.OrderItemID)
			if err != nil {
				s.logger.Error("[SEEDER] Error deleting order item", zap.Error(err), zap.Int64("order_item_id", int64(item.OrderItemID)))
				return err
			}
			deletedItemCount++
		}
	} else {
		s.logger.Error("[SEEDER] Failed to get order items list for clearing", zap.Error(err))
		return err
	}

	deletedOrderCount := 0
	orders, err := s.db.GetOrders(ctx, db.GetOrdersParams{Limit: 1000, Offset: 0})
	if err == nil {
		for _, order := range orders {
			err = s.db.DeleteOrderPermanently(ctx, order.OrderID)
			if err != nil {
				s.logger.Error("[SEEDER] Error deleting order", zap.Error(err), zap.Int64("order_id", int64(order.OrderID)))
				return err
			}
			deletedOrderCount++
		}
	} else {
		s.logger.Error("[SEEDER] Failed to get orders list for clearing", zap.Error(err))
		return err
	}

	s.logger.Info("[SEEDER] Orders service data cleared successfully",
		zap.Int("deleted_order_items", deletedItemCount),
		zap.Int("deleted_orders", deletedOrderCount),
	)
	return nil
}
