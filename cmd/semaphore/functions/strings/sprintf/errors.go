package sprintf

import (
	"errors"
	"fmt"
)

var (
	errUnknownFormatter   = errors.New("unable to detect the formatter")
	errIncomplete         = errors.New("verb is missing")
	errMalformedPrecision = errors.New("malformed precision")
	errTypeMismatch       = errors.New("type of argument does not match the verb")
	errNoValue            = errors.New("value is not set")
	errNonIntegerType     = errors.New("not an integer")
	errNonStringType      = errors.New("not a string")
	errNonFloatType       = errors.New("not a float")
)

type errVerbConflict struct {
	fmt.Stringer
}

func (e errVerbConflict) Error() string {
	return fmt.Sprintf("verb %q is already in use", e.Stringer)
}
