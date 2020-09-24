package template

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/prettyerr"
)

// ErrPathNotFound occurs when path cannot be resolved
type ErrPathNotFound struct {
	Path string
}

// Error returns a description of the given error as a string
func (e ErrPathNotFound) Error() string {
	return fmt.Sprintf("unable to resolve path '%s'", e.Path)
}

// Prettify returns the prettified version of the given error
func (e ErrPathNotFound) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "PathNotFound",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Path": e.Path,
		},
	}
}
