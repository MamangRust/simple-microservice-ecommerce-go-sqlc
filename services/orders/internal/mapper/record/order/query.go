package orderrecordmapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/record"
	db "github.com/MamangRust/simple_microservice_ecommerce/order/pkg/database/schema"
)

type orderQueryRecordMapper struct {
}

func NewOrderQueryRecordMapper() OrderQueryRecordMapper {
	return &orderQueryRecordMapper{}
}

func (s *orderQueryRecordMapper) ToOrderRecord(order *db.Order) *record.OrderRecord {
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

func (s *orderQueryRecordMapper) ToOrdersRecord(orders []*db.Order) []*record.OrderRecord {
	var result []*record.OrderRecord

	for _, order := range orders {
		result = append(result, s.ToOrderRecord(order))
	}

	return result
}

func (s *orderQueryRecordMapper) ToOrderRecordPagination(order *db.GetOrdersRow) *record.OrderRecord {
	var deletedAt *string
	if order.DeletedAt.Valid {
		deletedAtStr := order.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.OrderRecord{
		ID:         int(order.OrderID),
		TotalPrice: int(order.TotalPrice),
		CreatedAt:  order.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:  order.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:  deletedAt,
	}
}

func (s *orderQueryRecordMapper) ToOrdersRecordPagination(orders []*db.GetOrdersRow) []*record.OrderRecord {
	var result []*record.OrderRecord

	for _, order := range orders {
		result = append(result, s.ToOrderRecordPagination(order))
	}

	return result
}

func (s *orderQueryRecordMapper) ToOrderRecordActivePagination(order *db.GetOrdersActiveRow) *record.OrderRecord {
	var deletedAt *string
	if order.DeletedAt.Valid {
		deletedAtStr := order.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.OrderRecord{
		ID:         int(order.OrderID),
		TotalPrice: int(order.TotalPrice),
		CreatedAt:  order.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:  order.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:  deletedAt,
	}
}

func (s *orderQueryRecordMapper) ToOrdersRecordActivePagination(orders []*db.GetOrdersActiveRow) []*record.OrderRecord {
	var result []*record.OrderRecord

	for _, order := range orders {
		result = append(result, s.ToOrderRecordActivePagination(order))
	}

	return result
}

func (s *orderQueryRecordMapper) ToOrderRecordTrashedPagination(order *db.GetOrdersTrashedRow) *record.OrderRecord {
	var deletedAt *string
	if order.DeletedAt.Valid {
		deletedAtStr := order.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.OrderRecord{
		ID:         int(order.OrderID),
		TotalPrice: int(order.TotalPrice),
		CreatedAt:  order.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:  order.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:  deletedAt,
	}
}

func (s *orderQueryRecordMapper) ToOrdersRecordTrashedPagination(orders []*db.GetOrdersTrashedRow) []*record.OrderRecord {
	var result []*record.OrderRecord

	for _, order := range orders {
		result = append(result, s.ToOrderRecordTrashedPagination(order))
	}

	return result
}
