package orderitemrecordmapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/order/pkg/database/schema"
)

type orderItemQueryRecordMapper struct{}

func NewOrderItemQueryRecordMapper() OrderItemQueryRecordMapper {
	return &orderItemQueryRecordMapper{}
}

func (s *orderItemQueryRecordMapper) ToOrderItemRecord(orderItems *db.OrderItem) *record.OrderItemRecord {
	var deletedAt *string
	if orderItems.DeletedAt.Valid {
		deletedAtStr := orderItems.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.OrderItemRecord{
		ID:        int(orderItems.OrderItemID),
		OrderID:   int(orderItems.OrderID),
		ProductID: int(orderItems.ProductID),
		Quantity:  int(orderItems.Quantity),
		Price:     int(orderItems.Price),
		CreatedAt: orderItems.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: orderItems.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: deletedAt,
	}
}

func (s *orderItemQueryRecordMapper) ToOrderItemsRecord(orders []*db.OrderItem) []*record.OrderItemRecord {
	var result []*record.OrderItemRecord

	for _, Merchant := range orders {
		result = append(result, s.ToOrderItemRecord(Merchant))
	}

	return result
}

func (s *orderItemQueryRecordMapper) ToOrderItemRecordPagination(orderItems *db.GetOrderItemsRow) *record.OrderItemRecord {
	var deletedAt *string
	if orderItems.DeletedAt.Valid {
		deletedAtStr := orderItems.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.OrderItemRecord{
		ID:        int(orderItems.OrderItemID),
		OrderID:   int(orderItems.OrderID),
		ProductID: int(orderItems.ProductID),
		Quantity:  int(orderItems.Quantity),
		Price:     int(orderItems.Price),
		CreatedAt: orderItems.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: orderItems.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: deletedAt,
	}
}

func (s *orderItemQueryRecordMapper) ToOrderItemsRecordPagination(orders []*db.GetOrderItemsRow) []*record.OrderItemRecord {
	var result []*record.OrderItemRecord

	for _, Merchant := range orders {
		result = append(result, s.ToOrderItemRecordPagination(Merchant))
	}

	return result
}

func (s *orderItemQueryRecordMapper) ToOrderItemRecordActivePagination(orderItems *db.GetOrderItemsActiveRow) *record.OrderItemRecord {
	var deletedAt *string
	if orderItems.DeletedAt.Valid {
		deletedAtStr := orderItems.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.OrderItemRecord{
		ID:        int(orderItems.OrderItemID),
		OrderID:   int(orderItems.OrderID),
		ProductID: int(orderItems.ProductID),
		Quantity:  int(orderItems.Quantity),
		Price:     int(orderItems.Price),
		CreatedAt: orderItems.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: orderItems.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: deletedAt,
	}
}

func (s *orderItemQueryRecordMapper) ToOrderItemsRecordActivePagination(orders []*db.GetOrderItemsActiveRow) []*record.OrderItemRecord {
	var result []*record.OrderItemRecord

	for _, Merchant := range orders {
		result = append(result, s.ToOrderItemRecordActivePagination(Merchant))
	}

	return result
}

func (s *orderItemQueryRecordMapper) ToOrderItemRecordTrashedPagination(orderItems *db.GetOrderItemsTrashedRow) *record.OrderItemRecord {
	var deletedAt *string
	if orderItems.DeletedAt.Valid {
		deletedAtStr := orderItems.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.OrderItemRecord{
		ID:        int(orderItems.OrderItemID),
		OrderID:   int(orderItems.OrderID),
		ProductID: int(orderItems.ProductID),
		Quantity:  int(orderItems.Quantity),
		Price:     int(orderItems.Price),
		CreatedAt: orderItems.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: orderItems.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: deletedAt,
	}
}

func (s *orderItemQueryRecordMapper) ToOrderItemsRecordTrashedPagination(orders []*db.GetOrderItemsTrashedRow) []*record.OrderItemRecord {
	var result []*record.OrderItemRecord

	for _, Merchant := range orders {
		result = append(result, s.ToOrderItemRecordTrashedPagination(Merchant))
	}

	return result
}
