package orderitemrecordmapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/order/pkg/database/schema"
)

type orderItemCommandRecordMapper struct{}

func NewOrderItemCommandRecordMapper() OrderItemCommandRecordMapper {
	return &orderItemCommandRecordMapper{}
}

func (s *orderItemCommandRecordMapper) ToOrderItemRecord(orderItems *db.OrderItem) *record.OrderItemRecord {
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
