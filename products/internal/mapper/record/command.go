package productrecordmapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/product/pkg/database/schema"
)

type productCommandMapper struct {
}

func NewProductCommandRecordMapper() ProductCommandRecordMapping {
	return &productCommandMapper{}
}

func (s *productCommandMapper) ToProductRecord(product *db.Product) *record.ProductRecord {
	deletedAt := product.DeletedAt.Time.Format("2006-01-02 15:04:05.000")

	return &record.ProductRecord{
		ID:        int(product.ProductID),
		Name:      product.Name,
		Price:     int(product.Price),
		Stock:     int(product.Stock),
		CreatedAt: product.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: product.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: &deletedAt,
	}
}
