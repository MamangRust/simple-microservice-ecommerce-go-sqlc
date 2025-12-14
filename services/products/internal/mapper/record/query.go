package productrecordmapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/product/pkg/database/schema"
)

type productQueryMapper struct {
}

func NewProductQueryRecordMapper() ProductQueryRecordMapping {
	return &productQueryMapper{}
}

func (s *productQueryMapper) ToProductRecord(product *db.Product) *record.ProductRecord {
	var deletedAt *string
	if product.DeletedAt.Valid {
		formatted := product.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &formatted
	}

	return &record.ProductRecord{
		ID:        int(product.ProductID),
		Name:      product.Name,
		Price:     int(product.Price),
		Stock:     int(product.Stock),
		CreatedAt: product.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: product.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: deletedAt,
	}
}

func (s *productQueryMapper) ToProductsRecord(products []*db.Product) []*record.ProductRecord {
	var result []*record.ProductRecord
	for _, product := range products {
		result = append(result, s.ToProductRecord(product))
	}
	return result
}

func (s *productQueryMapper) ToProductRecordAll(product *db.GetProductsRow) *record.ProductRecord {
	var deletedAt *string
	if product.DeletedAt.Valid {
		formatted := product.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &formatted
	}

	return &record.ProductRecord{
		ID:        int(product.ProductID),
		Name:      product.Name,
		Price:     int(product.Price),
		Stock:     int(product.Stock),
		CreatedAt: product.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: product.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: deletedAt,
	}
}

func (s *productQueryMapper) ToProductsRecordAll(products []*db.GetProductsRow) []*record.ProductRecord {
	var result []*record.ProductRecord
	for _, product := range products {
		result = append(result, s.ToProductRecordAll(product))
	}
	return result
}

func (s *productQueryMapper) ToProductRecordActive(product *db.GetActiveProductsRow) *record.ProductRecord {
	var deletedAt *string
	if product.DeletedAt.Valid {
		formatted := product.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &formatted
	}

	return &record.ProductRecord{
		ID:        int(product.ProductID),
		Name:      product.Name,
		Price:     int(product.Price),
		Stock:     int(product.Stock),
		CreatedAt: product.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: product.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: deletedAt,
	}
}

func (s *productQueryMapper) ToProductsRecordActive(products []*db.GetActiveProductsRow) []*record.ProductRecord {
	var result []*record.ProductRecord
	for _, product := range products {
		result = append(result, s.ToProductRecordActive(product))
	}
	return result
}

func (s *productQueryMapper) ToProductRecordTrashed(product *db.GetTrashedProductsRow) *record.ProductRecord {
	var deletedAt *string
	if product.DeletedAt.Valid {
		formatted := product.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &formatted
	}

	return &record.ProductRecord{
		ID:        int(product.ProductID),
		Name:      product.Name,
		Price:     int(product.Price),
		Stock:     int(product.Stock),
		CreatedAt: product.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: product.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: deletedAt,
	}
}

func (s *productQueryMapper) ToProductsRecordTrashed(products []*db.GetTrashedProductsRow) []*record.ProductRecord {
	var result []*record.ProductRecord
	for _, product := range products {
		result = append(result, s.ToProductRecordTrashed(product))
	}
	return result
}
