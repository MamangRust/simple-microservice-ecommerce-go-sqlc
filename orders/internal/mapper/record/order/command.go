package orderrecordmapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/order/pkg/database/schema"
)

type orderCommandRecordMapper struct {
}

func NewOrderCommandRecordMapper() OrderCommandRecordMapper {
	return &orderCommandRecordMapper{}
}

func (s *orderCommandRecordMapper) ToOrderRecord(order *db.Order) *record.OrderRecord {
	var deletedAt *string
	if order.DeletedAt.Valid {
		deletedAtStr := order.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.OrderRecord{
		ID:         int(order.OrderID),
		TotalPrice: int(order.TotalPrice),
		UserID:     int(order.UserID),
		CreatedAt:  order.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:  order.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:  deletedAt,
	}
}
