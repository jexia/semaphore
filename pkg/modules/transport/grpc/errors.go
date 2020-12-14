package grpc

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/prettyerr"
)

// ErrUnknownMethod occurs when undefined method is called
type ErrUnknownMethod struct {
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
