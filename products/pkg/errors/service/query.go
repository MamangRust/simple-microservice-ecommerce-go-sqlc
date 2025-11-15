package productserviceerror

import (
	"net/http"

	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/response"
)

var (
	ErrFailedCountStock = response.NewErrorResponse("Failed to count stock", http.StatusInternalServerError)

	ErrFailedDeletingNotFoundProduct = response.NewErrorResponse("Product not found", http.StatusNotFound)
	ErrFailedDeleteImageProduct      = response.NewErrorResponse("Failed to delete image product", http.StatusInternalServerError)

	ErrFailedFindAllProducts        = response.NewErrorResponse("Failed to find all products", http.StatusInternalServerError)
	ErrFailedFindProductsByMerchant = response.NewErrorResponse("Failed to find products by merchant", http.StatusInternalServerError)
	ErrFailedFindProductsByCategory = response.NewErrorResponse("Failed to find products by category", http.StatusInternalServerError)
	ErrFailedFindProductById        = response.NewErrorResponse("Failed to find product by ID", http.StatusInternalServerError)
	ErrFailedFindProductTrashedById = response.NewErrorResponse("Failed to find trashed product by ID", http.StatusInternalServerError)
	ErrFailedFindProductByTrashed   = response.NewErrorResponse("Failed to find product by trashed", http.StatusInternalServerError)

	ErrFailedFindProductsByActive  = response.NewErrorResponse("Failed to find active products", http.StatusInternalServerError)
	ErrFailedFindProductsByTrashed = response.NewErrorResponse("Failed to find trashed products", http.StatusInternalServerError)
)
