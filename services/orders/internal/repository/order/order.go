package orderrepository

import (
	orderrecordmapper "github.com/MamangRust/simple_microservice_ecommerce/order/internal/mapper/record/order"
	db "github.com/MamangRust/simple_microservice_ecommerce/order/pkg/database/schema"
)

type OrderRepository interface {
	OrderQueryRepo() OrderQueryRepository
	OrderCommandRepo() OrderCommandRepository
}

type orderRepository struct {
	orderQueryRepository   OrderQueryRepository
	orderCommandRepository OrderCommandRepository
}

func (r *orderRepository) OrderQueryRepo() OrderQueryRepository {
	return r.orderQueryRepository
}

func (r *orderRepository) OrderCommandRepo() OrderCommandRepository {
	return r.orderCommandRepository
}

func NewOrderRepository(db *db.Queries) OrderRepository {
	orderQueryMapper := orderrecordmapper.NewOrderQueryRecordMapper()
	orderCommandMapper := orderrecordmapper.NewOrderCommandRecordMapper()

	return &orderRepository{
		orderQueryRepository:   NewOrderQueryRepository(db, orderQueryMapper),
		orderCommandRepository: NewOrderCommandRepository(db, orderCommandMapper),
	}
}
