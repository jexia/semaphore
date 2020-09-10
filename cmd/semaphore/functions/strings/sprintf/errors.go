package sprintf

import (
	"errors"
	"fmt"
)

var (
	errUnknownFormatter   = fmt.Errorf("unknown formatter")
	errMissingFormat      = errors.New("format is missing")
	errMalformedPrecision = errors.New("malformed precision")
	errTypeMismatch       = errors.New("type of argument does not match the verb")
	errNoValue            = errors.New("value is not set")
)
