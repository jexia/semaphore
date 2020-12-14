package json

import "github.com/jexia/semaphore/pkg/prettyerr"

// ErrUndefinedSpecs occurs when spacs are nil
type ErrUndefinedSpecs struct{}

// Error returns a description of the given error as a string
func (e ErrUndefinedSpecs) Error() string {
	return "no object specs defined"
}

// Prettify returns the prettified version of the given error
func (e ErrUndefinedSpecs) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UndefinedSpecs",
		Message: e.Error(),
	}
}
