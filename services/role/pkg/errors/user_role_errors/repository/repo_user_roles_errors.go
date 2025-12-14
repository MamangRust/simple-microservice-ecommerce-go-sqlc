package userrolerepositoryerrors

import "errors"

var ErrAssignRoleToUser = errors.New("failed to assign role to user")

var ErrUpdateRoleToUser = errors.New("failed to assign role to user")

var ErrRemoveRole = errors.New("failed to remove role from user")
