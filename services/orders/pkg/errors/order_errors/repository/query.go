package orderrepositoryerrors

import "errors"

var (
	ErrFindAllOrders  = errors.New("failed to find all orders")
	ErrFindByActive   = errors.New("failed to find active orders")
	ErrFindByTrashed  = errors.New("failed to find trashed orders")
	ErrFindByMerchant = errors.New("failed to find orders by merchant")
	ErrFindById       = errors.New("failed to find order by ID")
)
