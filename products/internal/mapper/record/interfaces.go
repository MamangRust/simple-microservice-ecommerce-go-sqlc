package productrecordmapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/product/pkg/database/schema"
)

type ProductBaseRecordMapper interface {
	ToProductRecord(Product *db.Product) *record.ProductRecord
}

type ProductQueryRecordMapping interface {
	ProductBaseRecordMapper
	ToProductsRecord(Products []*db.Product) []*record.ProductRecord
	ToProductRecordAll(Product *db.GetProductsRow) *record.ProductRecord
	ToProductsRecordAll(Products []*db.GetProductsRow) []*record.ProductRecord
	ToProductRecordActive(Product *db.GetActiveProductsRow) *record.ProductRecord
	ToProductsRecordActive(Products []*db.GetActiveProductsRow) []*record.ProductRecord

	ToProductRecordTrashed(Product *db.GetTrashedProductsRow) *record.ProductRecord

	ToProductsRecordTrashed(Products []*db.GetTrashedProductsRow) []*record.ProductRecord
}

type ProductCommandRecordMapping interface {
	ProductBaseRecordMapper
}
