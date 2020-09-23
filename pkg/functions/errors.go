package functions

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/prettyerr"
)

// ErrUndefinedFunction occurs when custom function is not defined in property but referenced
type ErrUndefinedFunction struct {
	Function string
	Property string
}

// Error returns a description of the given error as a string
func (e ErrUndefinedFunction) Error() string {
	return fmt.Sprintf("undefined custom function '%s' in '%s'", e.Function, e.Property)
}

// Prettify returns the prettified version of the given error
func (e ErrUndefinedFunction) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UndefinedFunction",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Function": e.Function,
			"Property": e.Property,
		},
	}
}
