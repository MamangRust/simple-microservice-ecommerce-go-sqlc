package orderitemrecordmapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/order/pkg/database/schema"
)

type OrderItemBaseRecordMapper interface {
	ToOrderItemRecord(orderItems *db.OrderItem) *record.OrderItemRecord
}

type OrderItemQueryRecordMapper interface {
	OrderItemBaseRecordMapper

	ToOrderItemsRecord(orders []*db.OrderItem) []*record.OrderItemRecord

	ToOrderItemRecordPagination(OrderItem *db.GetOrderItemsRow) *record.OrderItemRecord
	ToOrderItemsRecordPagination(OrderItem []*db.GetOrderItemsRow) []*record.OrderItemRecord

	ToOrderItemRecordActivePagination(OrderItem *db.GetOrderItemsActiveRow) *record.OrderItemRecord
	ToOrderItemsRecordActivePagination(OrderItem []*db.GetOrderItemsActiveRow) []*record.OrderItemRecord
	ToOrderItemRecordTrashedPagination(OrderItem *db.GetOrderItemsTrashedRow) *record.OrderItemRecord
	ToOrderItemsRecordTrashedPagination(OrderItem []*db.GetOrderItemsTrashedRow) []*record.OrderItemRecord
}

type OrderItemCommandRecordMapper interface {
	OrderItemBaseRecordMapper
}
