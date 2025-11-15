package orderserviceerrors

import (
	"net/http"

	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
)

var (
	ErrFailedFindAllOrders        = response.NewErrorResponse("Failed to find all orders", http.StatusInternalServerError)
	ErrFailedFindOrderById        = response.NewErrorResponse("Failed to find order by ID", http.StatusInternalServerError)
	ErrFailedFindOrdersByActive   = response.NewErrorResponse("Failed to find active orders", http.StatusInternalServerError)
	ErrFailedFindOrdersByTrashed  = response.NewErrorResponse("Failed to find trashed orders", http.StatusInternalServerError)
	ErrFailedFindOrdersByMerchant = response.NewErrorResponse("Failed to find orders by merchant", http.StatusInternalServerError)
)
