package sprintf

import "errors"

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
