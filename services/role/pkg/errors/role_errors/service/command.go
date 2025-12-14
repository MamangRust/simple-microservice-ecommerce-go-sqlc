package roleserviceerrors

import (
	"net/http"

	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/response"
)

var ErrFailedCreateRole = response.NewErrorResponse("Failed to create Role", http.StatusInternalServerError)

var ErrFailedUpdateRole = response.NewErrorResponse("Failed to update Role", http.StatusInternalServerError)

var ErrFailedTrashedRole = response.NewErrorResponse("Failed to move Role to trash", http.StatusInternalServerError)

var ErrFailedRestoreRole = response.NewErrorResponse("Failed to restore Role", http.StatusInternalServerError)

var ErrFailedDeletePermanent = response.NewErrorResponse("Failed to delete Role permanently", http.StatusInternalServerError)

var ErrFailedRestoreAll = response.NewErrorResponse("Failed to restore all Roles", http.StatusInternalServerError)

var ErrFailedDeleteAll = response.NewErrorResponse("Failed to delete all Roles permanently", http.StatusInternalServerError)
