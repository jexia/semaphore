package sprintf

import (
	"errors"
	"fmt"
)

var (
	errUnknownFormatter   = fmt.Errorf("unknown formatter")
	errIncomplete         = errors.New("verb is missing")
	errMalformedPrecision = errors.New("malformed precision")
	errTypeMismatch       = errors.New("type of argument does not match the verb")
	errNoValue            = errors.New("value is not set")
)
