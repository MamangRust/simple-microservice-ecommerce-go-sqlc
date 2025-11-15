package refreshtokenrepositoryerror

import "errors"

// ErrParseDate is returned when parsing the expiration date of a token fails.
var ErrParseDate = errors.New("failed to parse expiration date")
