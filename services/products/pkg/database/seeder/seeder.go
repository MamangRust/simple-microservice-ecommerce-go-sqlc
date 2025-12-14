package seeder

import (
	"context"
	"fmt"
	"strings"

	db "github.com/MamangRust/simple_microservice_ecommerce/product/pkg/database/schema"
	"github.com/MamangRust/simple_microservice_ecommerce/product/pkg/logger"
	"go.uber.org/zap"
)

type ProductSeeder struct {
	db     *db.Queries
	logger logger.LoggerInterface
}

func NewProductSeeder(db *db.Queries, logger logger.LoggerInterface) *ProductSeeder {
	return &ProductSeeder{
		db:     db,
		logger: logger,
	}
}

func (s *ProductSeeder) SeedAll(ctx context.Context) error {
	s.logger.Info("[SEEDER] Starting products service seeding...")

	if err := s.SeedProducts(ctx); err != nil {
		s.logger.Error("[SEEDER] Failed to seed products", zap.Error(err))
		return fmt.Errorf("failed to seed products: %w", err)
	}

	s.logger.Info("[SEEDER] Products service seeding completed successfully")
	return nil
}

func (s *ProductSeeder) SeedProducts(ctx context.Context) error {
	s.logger.Info("[SEEDER] Seeding products...")

	products := []db.CreateProductParams{
		{Name: "Laptop Gaming Pro", Price: 15000000, Stock: 25},
		{Name: "Wireless Mouse", Price: 250000, Stock: 100},
		{Name: "Mechanical Keyboard", Price: 150000, Stock: 75},
		{Name: "USB-C Hub", Price: 250000, Stock: 50},
		{Name: "Monitor 27 inch", Price: 300000, Stock: 30},
		{Name: "Webcam HD", Price: 500000, Stock: 40},
		{Name: "Headphone Bluetooth", Price: 350000, Stock: 60},
		{Name: "Phone Stand", Price: 50000, Stock: 200},
		{Name: "Cable Management", Price: 75000, Stock: 150},
		{Name: "Desk Lamp LED", Price: 125000, Stock: 80},
	}

	createdCount := 0
	for _, product := range products {
		newProduct, err := s.db.CreateProduct(ctx, db.CreateProductParams{
			Name:  product.Name,
			Price: product.Price,
			Stock: product.Stock,
		})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "UNIQUE constraint failed") {
				s.logger.Info("[SEEDER] Product already exists, skipping", zap.String("product_name", product.Name))
				continue
			}
			s.logger.Error("[SEEDER] Error inserting product", zap.String("product_name", product.Name), zap.Error(err))
			return err
		}

		s.logger.Info("[SEEDER] Successfully created product",
			zap.String("product_name", product.Name),
			zap.Int64("product_id", int64(newProduct.ProductID)),
			zap.Int64("price", product.Price),
			zap.Int64("stock", int64(product.Stock)),
		)
		createdCount++
	}

	s.logger.Info("[SEEDER] Product seeding summary",
		zap.Int("total_products_in_list", len(products)),
		zap.Int("products_created", createdCount),
	)
	return nil
}

func (s *ProductSeeder) ClearAll(ctx context.Context) error {
	s.logger.Info("[SEEDER] Clearing products service data...")

	products, err := s.db.GetProducts(ctx, db.GetProductsParams{
		Column1: "",
		Limit:   1000,
		Offset:  0,
	})

	deletedCount := 0
	if err == nil {
		for _, product := range products {
			err = s.db.DeleteProductPermanently(ctx, product.ProductID)
			if err != nil {
				s.logger.Error("[SEEDER] Error deleting product",
					zap.Int64("product_id", int64(product.ProductID)),
					zap.String("product_name", product.Name),
					zap.Error(err),
				)
				return err
			}
			deletedCount++
		}
	} else {
		s.logger.Error("[SEEDER] Failed to get products list for clearing", zap.Error(err))
		return err
	}

	s.logger.Info("[SEEDER] Products service data cleared successfully",
		zap.Int("products_deleted", deletedCount),
	)
	return nil
}
