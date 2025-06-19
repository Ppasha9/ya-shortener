package errors

import (
	"errors"
)

var ErrOrigURLDuplicate = errors.New("ErrOrigURLDuplicate")
var ErrInvalidAuthCookie = errors.New("ErrInvalidAuthCookie")
var ErrInvalidAuthCookieBytesLen = errors.New("ErrInvalidAuthCookieBytesLen")
