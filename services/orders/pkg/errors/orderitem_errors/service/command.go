package orderitemserviceerrors

import (
	"net/http"

	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
)

var (
	ErrFailedRestoreAllOrderItem = response.NewErrorResponse("Failed to restore all order items", http.StatusInternalServerError)
	ErrFailedDeleteAllOrderItem  = response.NewErrorResponse("Failed to delete all order items", http.StatusInternalServerError)

	ErrFailedTrashedOrderItem = response.NewErrorResponse("Order item is already trashed", http.StatusBadRequest)
	ErrFailedRestoreOrderItem = response.NewErrorResponse("Failed to restore order item", http.StatusInternalServerError)
	ErrFailedDeleteOrderItem  = response.NewErrorResponse("Failed to delete order item", http.StatusInternalServerError)

	ErrFailedCreateOrderItem = response.NewErrorResponse("Failed to create order item", http.StatusInternalServerError)
	ErrFailedUpdateOrderItem = response.NewErrorResponse("Failed to update order item", http.StatusInternalServerError)
	ErrFailedCalculateTotal  = response.NewErrorResponse("Failed to calculate total", http.StatusInternalServerError)
)
