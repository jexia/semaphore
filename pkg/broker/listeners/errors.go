package listeners

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/prettyerr"
)

// ErrNoListener is thrown when no listener has been found
type ErrNoListener struct {
	Listener string
}

func (e ErrNoListener) Error() string {
	return fmt.Sprintf("unknown listener '%s'", e.Listener)
}

// Prettify returns the prettified version of the given error
func (e ErrNoListener) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Message: e.Error(),
		Details: map[string]interface{}{"listener": e.Listener},
	}
}
