package orderitemserviceerrors

import (
	"net/http"

	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
)

var (
	ErrFailedFindAllOrderItems       = response.NewErrorResponse("Failed to find all order items", http.StatusInternalServerError)
	ErrFailedFindOrderItemsByActive  = response.NewErrorResponse("Failed to find active order items", http.StatusInternalServerError)
	ErrFailedFindOrderItemsByTrashed = response.NewErrorResponse("Failed to find trashed order items", http.StatusInternalServerError)
	ErrFailedFindOrderItemByOrder    = response.NewErrorResponse("Failed to find order items by order ID", http.StatusInternalServerError)
)
