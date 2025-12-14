package orderrepositoryerrors

import "errors"

var (
	ErrCreateOrder             = errors.New("failed to create order")
	ErrUpdateOrder             = errors.New("failed to update order")
	ErrTrashedOrder            = errors.New("failed to move order to trash")
	ErrRestoreOrder            = errors.New("failed to restore order from trash")
	ErrDeleteOrderPermanent    = errors.New("failed to permanently delete order")
	ErrRestoreAllOrder         = errors.New("failed to restore all trashed orders")
	ErrDeleteAllOrderPermanent = errors.New("failed to permanently delete all trashed orders")
)
