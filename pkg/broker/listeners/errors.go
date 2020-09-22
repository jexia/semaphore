package listeners

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/prettyerr"
)

type ErrNoListener struct {
	Listener string
}

func (e ErrNoListener) Error() string {
	return fmt.Sprintf("unknown listener '%s'", e.Listener)
}

func (e ErrNoListener) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Message: e.Error(),
		Details: map[string]interface{}{"listener": e.Listener},
	}
}
