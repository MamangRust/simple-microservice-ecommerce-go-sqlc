package productserviceerror

import (
	"net/http"

	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/response"
)

var (
	ErrFailedCreateProduct = response.NewErrorResponse("Failed to create product", http.StatusInternalServerError)
	ErrFailedUpdateProduct = response.NewErrorResponse("Failed to update product", http.StatusInternalServerError)

	ErrFailedTrashProduct               = response.NewErrorResponse("Failed to trash product", http.StatusInternalServerError)
	ErrFailedRestoreProduct             = response.NewErrorResponse("Failed to restore product", http.StatusInternalServerError)
	ErrFailedDeleteProductPermanent     = response.NewErrorResponse("Failed to permanently delete product", http.StatusInternalServerError)
	ErrFailedRestoreAllProducts         = response.NewErrorResponse("Failed to restore all products", http.StatusInternalServerError)
	ErrFailedDeleteAllProductsPermanent = response.NewErrorResponse("Failed to permanently delete all products", http.StatusInternalServerError)
)
