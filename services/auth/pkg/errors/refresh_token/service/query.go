package refreshtokenserviceerrors

import (
	"net/http"

	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/response"
)

// ErrRefreshTokenNotFound indicates that the refresh token was not found.
var ErrRefreshTokenNotFound = response.NewErrorResponse("Refresh token not found", http.StatusNotFound)

// ErrFailedFindByToken indicates failure when searching for a refresh token by its token value.
var ErrFailedFindByToken = response.NewErrorResponse("Failed to find refresh token by token", http.StatusInternalServerError)

// ErrFailedFindByUserID indicates failure when searching for a refresh token by user ID.
var ErrFailedFindByUserID = response.NewErrorResponse("Failed to find refresh token by user ID", http.StatusInternalServerError)
