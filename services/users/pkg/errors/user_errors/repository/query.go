package userrepositoryerrors

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrFindAllUsers     = errors.New("failed to find all users")
	ErrFindActiveUsers  = errors.New("failed to find active users")
	ErrFindTrashedUsers = errors.New("failed to find trashed users")
)
