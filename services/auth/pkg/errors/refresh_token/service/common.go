package refreshtokenserviceerrors

import (
	"net/http"

	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/response"
)

// ErrFailedParseExpirationDate indicates failure when parsing the expiration date of a token.
var ErrFailedParseExpirationDate = response.NewErrorResponse("Failed to parse expiration date", http.StatusBadRequest)

// ErrFailedInValidToken is returned when an access token is invalid.
var ErrFailedInvalidToken = response.NewErrorResponse("Failed to invalid access token", http.StatusInternalServerError)

// ErrFailedInValidUserId is returned when a user ID is invalid.
var ErrFailedInvalidUserId = response.NewErrorResponse("Failed to invalid user id", http.StatusInternalServerError)

// ErrFailedExpire occurs when expiring a refresh token fails.
var ErrFailedExpire = response.NewErrorResponse("Failed to find refresh token by token", http.StatusInternalServerError)
