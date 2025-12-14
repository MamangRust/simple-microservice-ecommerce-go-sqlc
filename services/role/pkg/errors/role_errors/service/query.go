package roleserviceerrors

import (
	"net/http"

	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/response"
)

var ErrRoleNotFoundRes = response.NewErrorResponse("Role not found", http.StatusNotFound)

var ErrFailedFindAll = response.NewErrorResponse("Failed to fetch Roles", http.StatusInternalServerError)

var ErrFailedFindActive = response.NewErrorResponse("Failed to fetch active Roles", http.StatusInternalServerError)

var ErrFailedFindTrashed = response.NewErrorResponse("Failed to fetch trashed Roles", http.StatusInternalServerError)
