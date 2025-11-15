package userserviceerrors

import (
	"net/http"

	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/response"
)

var ErrInternalServerError = response.NewErrorResponse("Internal server error", http.StatusInternalServerError)

var ErrFailedSendEmail = response.NewErrorResponse("Failed to send email", http.StatusInternalServerError)

var ErrFailedPasswordNoMatch = response.NewErrorResponse("Failed password not match", http.StatusBadRequest)

var ErrUserPassword = response.NewErrorResponse("Failed invalid password", http.StatusBadRequest)
