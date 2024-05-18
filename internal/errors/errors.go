package customerrors

import (
	"errors"
)

// DuplicatedKey is an error that is returned when a key is duplicated
var DuplicatedKey = errors.New("duplicated key")
