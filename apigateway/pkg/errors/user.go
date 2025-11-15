package errors

import (
	"net/http"

	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/domain/response"
	"github.com/labstack/echo/v4"
)

var (
	ErrApiUserNotFound = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "User not found", http.StatusNotFound)
	}

	ErrApiUserInvalidId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid User id", http.StatusNotFound)
	}

	ErrApiUserFailedFindAll = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch Users", http.StatusInternalServerError)
	}

	ErrApiUserFailedFindActive = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch active Users", http.StatusInternalServerError)
	}

	ErrApiUserFailedFindTrashed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch trashed Users", http.StatusInternalServerError)
	}
	ErrInvalidUserId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid User id", http.StatusBadRequest)
	}

	ErrApiUserFailedCreateUser = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to create User", http.StatusInternalServerError)
	}

	ErrApiUserFailedUpdateUser = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to update User", http.StatusInternalServerError)
	}

	ErrApiUserValidateCreateUser = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid create User request", http.StatusBadRequest)
	}

	ErrApiUserValidateUpdateUser = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid update User request", http.StatusBadRequest)
	}

	ErrApiUserBindCreateUser = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid create User request", http.StatusBadRequest)
	}

	ErrApiUserBindUpdateUser = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid update User request", http.StatusBadRequest)
	}

	ErrApiUserFailedTrashedUser = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to move User to trash", http.StatusInternalServerError)
	}

	ErrApiUserFailedRestoreUser = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore User", http.StatusInternalServerError)
	}

	ErrApiUserFailedDeletePermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to delete User permanently", http.StatusInternalServerError)
	}

	ErrApiUserFailedRestoreAll = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore all Users", http.StatusInternalServerError)
	}

	ErrApiUserFailedDeleteAll = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to delete all Users permanently", http.StatusInternalServerError)
	}
)
