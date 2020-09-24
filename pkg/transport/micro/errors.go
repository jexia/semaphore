package micro

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/prettyerr"
)

type wrapErr struct {
	Inner error
}

func (i wrapErr) Unwrap() error {
	return i.Inner
}

// ErrUndefinedMethod occurs when undefined method is not defined
type ErrUndefinedMethod struct {
	wrapErr
}

// Error returns a description of the given error as a string
func (e ErrUndefinedMethod) Error() string {
	return fmt.Sprint("method required, proxy forward not supported")
}

// Prettify returns the prettified version of the given error
func (e ErrUndefinedMethod) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UndefinedMethod",
		Message: e.Error(),
		Details: map[string]interface{}{},
	}
}

// ErrUnknownMethod occurs when undefined method is called
type ErrUnknownMethod struct {
	wrapErr
	Method string
}

// Error returns a description of the given error as a string
func (e ErrUnknownMethod) Error() string {
	return fmt.Sprintf("unknown service method %s", e.Method)
}

// Prettify returns the prettified version of the given error
func (e ErrUnknownMethod) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UnknownMethod",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Method": e.Method,
		},
	}
}
