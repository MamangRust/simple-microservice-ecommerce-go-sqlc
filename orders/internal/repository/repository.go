package repository

import (
	orderrepository "github.com/MamangRust/simple_microservice_ecommerce/order/internal/repository/order"
	orderitemrepository "github.com/MamangRust/simple_microservice_ecommerce/order/internal/repository/order_item"
	db "github.com/MamangRust/simple_microservice_ecommerce/order/pkg/database/schema"
)

type Repository interface {
	orderrepository.OrderRepository
	orderitemrepository.OrderItemRepository
}

type repository struct {
	orderrepository.OrderRepository
	orderitemrepository.OrderItemRepository
}

func NewRepository(db *db.Queries) Repository {
	return &repository{
		OrderRepository:     orderrepository.NewOrderRepository(db),
		OrderItemRepository: orderitemrepository.NewOrderItemRepository(db),
	}
}
