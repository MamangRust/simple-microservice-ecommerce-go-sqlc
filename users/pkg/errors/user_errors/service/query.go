package userserviceerrors

import (
	"net/http"

	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/response"
)

var ErrUserNotFoundRes = response.NewErrorResponse("User not found", http.StatusNotFound)
var ErrUserEmailAlready = response.NewErrorResponse("User email already exists", http.StatusBadRequest)
var ErrFailedFindAll = response.NewErrorResponse("Failed to fetch users", http.StatusInternalServerError)
var ErrFailedFindActive = response.NewErrorResponse("Failed to fetch active users", http.StatusInternalServerError)
var ErrFailedFindTrashed = response.NewErrorResponse("Failed to fetch trashed users", http.StatusInternalServerError)
