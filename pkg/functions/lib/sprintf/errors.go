package sprintf

import (
	"errors"
	"fmt"

	"github.com/jexia/semaphore/pkg/specs"
)

var (
	errNoArguments        = errors.New("at least 1 argument is expected")
	errNoReferenceSupport = errors.New("cannot use reference as a format")
	errNoFormat           = errors.New("format is not set")
	errInvalidFormat      = errors.New("format must be a string")
	errUnknownFormatter   = errors.New("unable to detect the formatter")
	errIncomplete         = errors.New("verb is missing")
	errMalformedPrecision = errors.New("malformed precision")
	errTypeMismatch       = errors.New("type of argument does not match the verb")
	errNoValue            = errors.New("value is not set")
	errNonIntegerType     = errors.New("not an integer")
	errNonStringType      = errors.New("not a string")
	errNonFloatType       = errors.New("not a float")
)

type errInvalidArguments struct {
	actual, expected int
}

func (e errInvalidArguments) Error() string {
	return fmt.Sprintf("invalid number of input arguments %d, expected %d", e.actual, e.expected)
}

type errCannotFormat struct {
	formatter fmt.Stringer
	argument  *specs.Property
}

func (e errCannotFormat) Error() string {
	return fmt.Sprintf("cannot use '%%%s' formatter for argument '%s' of type '%s'", e.formatter, e.argument.Name, e.argument.Type())
}

type errVerbConflict struct {
	fmt.Stringer
}

func (e errVerbConflict) Error() string {
	return fmt.Sprintf("verb %q is already in use", e.Stringer)
}

type errScanFormat struct {
	inner    error
	format   string
	position int
}

func (e errScanFormat) Unwrap() error { return e.inner }

func (e errScanFormat) Error() string {
	var msg = e.inner.Error() + ":\n"
	for i := 0; i < e.position; i++ {
		msg += " "
	}

	return msg + "↓\n" + e.format + "\n"
}
