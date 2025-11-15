package orderrecordmapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/order/pkg/database/schema"
)

type OrderBaseRecordMapper interface {
	ToOrderRecord(order *db.Order) *record.OrderRecord
}

type OrderQueryRecordMapper interface {
	OrderBaseRecordMapper

	ToOrdersRecord(orders []*db.Order) []*record.OrderRecord
	ToOrderRecordPagination(order *db.GetOrdersRow) *record.OrderRecord
	ToOrdersRecordPagination(orders []*db.GetOrdersRow) []*record.OrderRecord
	ToOrderRecordActivePagination(order *db.GetOrdersActiveRow) *record.OrderRecord
	ToOrdersRecordActivePagination(orders []*db.GetOrdersActiveRow) []*record.OrderRecord
	ToOrderRecordTrashedPagination(order *db.GetOrdersTrashedRow) *record.OrderRecord
	ToOrdersRecordTrashedPagination(orders []*db.GetOrdersTrashedRow) []*record.OrderRecord
}

type OrderCommandRecordMapper interface {
	ToOrderRecord(order *db.Order) *record.OrderRecord
}
