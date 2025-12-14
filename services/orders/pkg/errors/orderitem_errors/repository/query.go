package orderitemrepositoryerrors

import "errors"

var (
	ErrFindAllOrderItems    = errors.New("failed to find all order items")
	ErrFindByActive         = errors.New("failed to find active order items")
	ErrFindByTrashed        = errors.New("failed to find trashed order items")
	ErrFindOrderItemByOrder = errors.New("failed to find order items by order ID")
)
