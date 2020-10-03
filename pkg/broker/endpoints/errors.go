package endpoints

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/prettyerr"
)

// ErrUnknownService occurs when there is no defined service for the requested method
type ErrUnknownService struct {
	Service string
}

// Prettify returns the prettified version of the given error
func (e ErrUnknownService) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Original: nil,
		Message:  e.Error(),
		Details: map[string]interface{}{
			"service": e.Service,
		},
	}
}

func (e ErrUnknownService) Error() string {
	return fmt.Sprintf("unknown service '%s'", e.Service)
}
