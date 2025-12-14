package userrepositoryerrors

import "errors"

var ErrUserConflict = errors.New("failed user already exists")

var ErrCreateUser = errors.New("failed to create user")

var ErrUpdateUser = errors.New("failed to update user")

var ErrUpdateUserVerificationCode = errors.New("failed to update user verification code")

var ErrUpdateUserPassword = errors.New("failed to update user password")

var ErrTrashedUser = errors.New("failed to move user to trash")

var ErrRestoreUser = errors.New("failed to restore user from trash")

var ErrDeleteUserPermanent = errors.New("failed to permanently delete user")

var ErrRestoreAllUsers = errors.New("failed to restore all users")

var ErrDeleteAllUsers = errors.New("failed to permanently delete all users")
