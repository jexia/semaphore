package formencoded

import (
	"errors"
	"fmt"

	"github.com/jexia/semaphore/pkg/prettyerr"
)

var errNilSchema = errors.New("nil schema")

type errUndefinedProperty string

func (e errUndefinedProperty) Error() string {
	return fmt.Sprintf("undefined property %q", string(e))
}

type errUnknownLabel string

func (e errUnknownLabel) Error() string {
	return fmt.Sprintf("unknown label %q", string(e))
}

// ErrUndefinedSpecs occurs when spacs are nil
type ErrUndefinedSpecs struct{}

// Error returns a description of the given error as a string
func (e ErrUndefinedSpecs) Error() string {
	return fmt.Sprint("no object specs defined")
}

// Prettify returns the prettified version of the given error
func (e ErrUndefinedSpecs) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UndefinedSpecs",
		Message: e.Error(),
	}
}
