package orderserviceerrors

import (
	"net/http"

	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
)

var (
	ErrFailedNotDeleteAtOrder = response.NewErrorResponse("Failed to delete at order", http.StatusInternalServerError)

	ErrInsufficientProductStock  = response.NewErrorResponse("Insufficient product stock", http.StatusBadRequest)
	ErrFailedInvalidCountInStock = response.NewErrorResponse("Failed to find invalid count in stock", http.StatusInternalServerError)

	ErrFailedCreateOrder             = response.NewErrorResponse("Failed to create order", http.StatusInternalServerError)
	ErrFailedUpdateOrder             = response.NewErrorResponse("Failed to update order", http.StatusInternalServerError)
	ErrFailedTrashOrder              = response.NewErrorResponse("Failed to trash order", http.StatusInternalServerError)
	ErrFailedRestoreOrder            = response.NewErrorResponse("Failed to restore order", http.StatusInternalServerError)
	ErrFailedDeleteOrderPermanent    = response.NewErrorResponse("Failed to permanently delete order", http.StatusInternalServerError)
	ErrFailedRestoreAllOrder         = response.NewErrorResponse("Failed to restore all orders", http.StatusInternalServerError)
	ErrFailedDeleteAllOrderPermanent = response.NewErrorResponse("Failed to permanently delete all orders", http.StatusInternalServerError)
)
