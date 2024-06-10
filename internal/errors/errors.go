package customerrors

import (
	"errors"
)

var (
	ErrDuplicatedKey   = errors.New("duplicated key")
	ErrInvalidTimezone = errors.New("invalid timezone")
)
