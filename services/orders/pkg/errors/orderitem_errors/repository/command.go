package orderitemrepositoryerrors

import "errors"

var (
	ErrCalculateTotalPrice      = errors.New("failed to calculate total price")
	ErrCreateOrderItem          = errors.New("failed to create order item")
	ErrUpdateOrderItem          = errors.New("failed to update order item")
	ErrTrashedOrderItem         = errors.New("failed to move order item to trash")
	ErrRestoreOrderItem         = errors.New("failed to restore order item from trash")
	ErrDeleteOrderItemPermanent = errors.New("failed to permanently delete order item")
	ErrRestoreAllOrderItem      = errors.New("failed to restore all trashed order items")
	ErrDeleteAllOrderPermanent  = errors.New("failed to permanently delete all trashed order items")
)
