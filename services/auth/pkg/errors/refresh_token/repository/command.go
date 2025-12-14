package refreshtokenrepositoryerror

import "errors"

// ErrCreateRefreshToken indicates that an error occurred while creating a new refresh token.
var ErrCreateRefreshToken = errors.New("failed to create refresh token")

// ErrUpdateRefreshToken is returned when the refresh token update process fails.
var ErrUpdateRefreshToken = errors.New("failed to update refresh token")

// ErrDeleteRefreshToken is returned when deleting a refresh token fails.
var ErrDeleteRefreshToken = errors.New("failed to delete refresh token")

// ErrDeleteByUserID indicates a failure when attempting to delete a refresh token using the user ID.
var ErrDeleteByUserID = errors.New("failed to delete refresh token by user ID")
