package rolerepositoryerrors

import "errors"

var ErrRoleNotFound = errors.New("role not found")

var ErrFindAllRoles = errors.New("failed to find all Roles")

var ErrFindActiveRoles = errors.New("failed to find active Roles")

var ErrFindTrashedRoles = errors.New("failed to find trashed Roles")

var ErrRoleConflict = errors.New("failed Role already exists")
