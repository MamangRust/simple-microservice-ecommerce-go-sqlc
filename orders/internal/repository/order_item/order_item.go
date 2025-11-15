package orderitemrepository

import (
	orderitemrecordmapper "github.com/MamangRust/simple_microservice_ecommerce/order/internal/mapper/record/orderitem"
	db "github.com/MamangRust/simple_microservice_ecommerce/order/pkg/database/schema"
)

type OrderItemRepository interface {
	OrderItemQueryRepo() OrderItemQueryRepository
	OrderItemCommandRepo() OrderItemCommandRepository
}

type orderItemRepository struct {
	orderItemQueryRepository   OrderItemQueryRepository
	orderItemCommandRepository OrderItemCommandRepository
}

func (r *orderItemRepository) OrderItemQueryRepo() OrderItemQueryRepository {
	return r.orderItemQueryRepository
}

func (r *orderItemRepository) OrderItemCommandRepo() OrderItemCommandRepository {
	return r.orderItemCommandRepository
}

func NewOrderItemRepository(db *db.Queries) OrderItemRepository {
	orderQueryMapper := orderitemrecordmapper.NewOrderItemQueryRecordMapper()
	orderCommandMapper := orderitemrecordmapper.NewOrderItemCommandRecordMapper()

	return &orderItemRepository{
		orderItemQueryRepository:   NewOrderItemQueryRepository(db, orderQueryMapper),
		orderItemCommandRepository: NewOrderItemCommandRepository(db, orderCommandMapper),
	}
}
