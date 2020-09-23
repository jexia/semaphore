package endpoints

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/prettyerr"
)

// ErrNoServiceForMethod occurs when there is no defined service for the requested method
type ErrNoServiceForMethod struct {
	Method string
}

func (e ErrNoServiceForMethod) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Original: nil,
		Message:  e.Error(),
		Details: map[string]interface{}{
			"method": e.Method,
		},
	}
}

func (e ErrNoServiceForMethod) Error() string {
	return fmt.Sprintf("failed to find service for '%s'", e.Method)
}
