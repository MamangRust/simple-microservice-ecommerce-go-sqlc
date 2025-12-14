package userserviceerrors

import (
	"net/http"

	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/response"
)

var ErrUserNotFoundRes = response.NewErrorResponse("User not found", http.StatusNotFound)

// ErrInternalServerError is a generic internal server error.
var ErrInternalServerError = response.NewErrorResponse("Internal server error", http.StatusInternalServerError)

// ErrFailedSendEmail is returned when sending an email fails.
var ErrFailedSendEmail = response.NewErrorResponse("Failed to send email", http.StatusInternalServerError)

// ErrFailedPasswordNoMatch is returned when passwords do not match.
var ErrFailedPasswordNoMatch = response.NewErrorResponse("Failed password not match", http.StatusBadRequest)

// ErrUserPassword is returned when there is an invalid password.
var ErrUserPassword = response.NewErrorResponse("Failed invalid password", http.StatusBadRequest)
