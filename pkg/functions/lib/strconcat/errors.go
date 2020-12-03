package strconcat

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/prettyerr"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

type wrapErr struct {
	Inner error
}

func (i wrapErr) Unwrap() error {
	return i.Inner
}

// ErrInvalidArgument is thrown when a given argument is invalid
type ErrInvalidArgument struct {
	wrapErr
	Function string
	Property *specs.Property
	Expected types.Type
}

// Error returns a description of the given error as a string
func (e ErrInvalidArgument) Error() string {
	return fmt.Sprintf("invalid argument %s in %s expected %s", e.Property.Type(), e.Function, e.Expected)
}

// Prettify returns the prettified version of the given error
func (e ErrInvalidArgument) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "InvalidArgument",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Function": e.Function,
			"Type":     e.Property.Type,
			"Expected": e.Expected,
		},
	}
}
