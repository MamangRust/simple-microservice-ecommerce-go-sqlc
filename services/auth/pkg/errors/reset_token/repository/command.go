package resettokenrepositoryerror

import "errors"

var ErrCreateResetToken = errors.New("failed to create reset token")

var ErrDeleteResetToken = errors.New("failed to delete reset token")
