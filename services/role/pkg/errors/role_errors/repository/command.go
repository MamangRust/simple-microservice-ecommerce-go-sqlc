package rolerepositoryerrors

import "errors"

var ErrCreateRole = errors.New("failed to create Role")

var ErrUpdateRole = errors.New("failed to update Role")

var ErrTrashedRole = errors.New("failed to move Role to trash")

var ErrRestoreRole = errors.New("failed to restore Role from trash")

var ErrDeleteRolePermanent = errors.New("failed to permanently delete Role")

var ErrRestoreAllRoles = errors.New("failed to restore all Roles")

var ErrDeleteAllRoles = errors.New("failed to permanently delete all Roles")
