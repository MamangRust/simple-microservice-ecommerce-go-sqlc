package refreshtokenserviceerrors

import (
	"net/http"

	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/response"
)

// ErrFailedCreateAccess is returned when the creation of an access token fails.
var ErrFailedCreateAccess = response.NewErrorResponse("Failed to create access token", http.StatusInternalServerError)

// ErrFailedCreateRefresh is returned when the creation of a refresh token fails.
var ErrFailedCreateRefresh = response.NewErrorResponse("Failed to create refresh token", http.StatusInternalServerError)

// ErrFailedCreateRefreshToken is returned when refresh token creation fails.
var ErrFailedCreateRefreshToken = response.NewErrorResponse("Failed to create refresh token", http.StatusInternalServerError)

// ErrFailedUpdateRefreshToken is returned when refresh token update fails.
var ErrFailedUpdateRefreshToken = response.NewErrorResponse("Failed to update refresh token", http.StatusInternalServerError)

// ErrFailedDeleteRefreshToken is returned when refresh token deletion fails.
var ErrFailedDeleteRefreshToken = response.NewErrorResponse("Failed to delete refresh token", http.StatusInternalServerError)

// ErrFailedDeleteByUserID is returned when deletion of a refresh token by user ID fails.
var ErrFailedDeleteByUserID = response.NewErrorResponse("Failed to delete refresh token by user ID", http.StatusInternalServerError)
