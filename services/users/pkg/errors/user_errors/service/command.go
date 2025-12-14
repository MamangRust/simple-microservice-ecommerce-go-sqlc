package userserviceerrors

import (
	"net/http"

	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/domain/response"
)

var ErrFailedCreateUser = response.NewErrorResponse("Failed to create user", http.StatusInternalServerError)

var ErrFailedUpdateUser = response.NewErrorResponse("Failed to update user", http.StatusInternalServerError)

var ErrFailedTrashedUser = response.NewErrorResponse("Failed to move user to trash", http.StatusInternalServerError)

var ErrFailedRestoreUser = response.NewErrorResponse("Failed to restore user", http.StatusInternalServerError)

var ErrFailedDeletePermanent = response.NewErrorResponse("Failed to delete user permanently", http.StatusInternalServerError)

var ErrFailedRestoreAll = response.NewErrorResponse("Failed to restore all users", http.StatusInternalServerError)

var ErrFailedDeleteAll = response.NewErrorResponse("Failed to delete all users permanently", http.StatusInternalServerError)
