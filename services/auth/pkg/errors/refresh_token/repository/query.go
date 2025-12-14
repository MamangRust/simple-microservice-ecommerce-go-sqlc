package refreshtokenrepositoryerror

import "errors"

// ErrTokenNotFound indicates that the refresh token could not be found.
var ErrTokenNotFound = errors.New("refresh token not found")

// ErrFindByToken is returned when a lookup for the refresh token by token value fails.
var ErrFindByToken = errors.New("failed to find refresh token by token")

// ErrFindByUserID is returned when a lookup for the refresh token by user ID fails.
var ErrFindByUserID = errors.New("failed to find refresh token by user ID")
