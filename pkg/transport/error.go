package transport

import (
	"github.com/jexia/maestro/pkg/specs"
)

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an transport Error returning error.
// Otherwise, Unwrap returns nil.
func Unwrap(err error) Error {
	u, ok := err.(Error)
	if !ok {
		return nil
	}

	return u
}

// WrapError returns an error with the error handle specs
func WrapError(err error, handle specs.ErrorHandle) Error {
	return &wrapper{
		err:         err,
		ErrorHandle: handle,
	}
}

// Error represents a wrapped error and error specs
type Error interface {
	specs.ErrorHandle
	String() string
	Error() string
}

type wrapper struct {
	specs.ErrorHandle
	err error
}

func (w *wrapper) String() string {
	if w.err == nil {
		return ""
	}

	return w.err.Error()
}

// Error returns the underlaying error as a string
func (w *wrapper) Error() string {
	if w.err == nil {
		return ""
	}

	return w.err.Error()
}

// Unwrap unwraps the given error and returns the wrapped error
func (w *wrapper) Unwrap() error {
	return w.err
}
