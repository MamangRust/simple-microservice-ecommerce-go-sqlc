package roleserviceerrors

import (
	"net/http"

	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/response"
)

var ErrRoleNotFoundRes = response.NewErrorResponse("Role not found", http.StatusNotFound)
