package proto

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/prettyerr"
)

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

// ErrNonRootMessage occurs when message type is not root
type ErrNonRootMessage struct{}

// Error returns a description of the given error as a string
func (e ErrNonRootMessage) Error() string {
	return fmt.Sprint("protobuffer messages root property should be a message")
}

// Prettify returns the prettified version of the given error
func (e ErrNonRootMessage) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "NonRootMessage",
		Message: e.Error(),
	}
}
