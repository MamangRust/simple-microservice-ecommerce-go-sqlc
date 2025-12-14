package userroleserviceerrors

import (
	"net/http"

	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/domain/response"
)

var ErrFailedAssignRoleToUser = response.NewErrorResponse(
	"Failed to assign role to user",
	http.StatusInternalServerError,
)

var ErrFailedUpdateRoleToUser = response.NewErrorResponse(
	"Failed to update role to user",
	http.StatusInternalServerError,
)

var ErrFailedRemoveRole = response.NewErrorResponse(
	"Failed to remove role from user",
	http.StatusInternalServerError,
)
