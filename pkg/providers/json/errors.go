package json

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/prettyerr"
)

type wrapErr struct {
	Inner error
}

func (i wrapErr) Unwrap() error {
	return i.Inner
}

// ErrPathNotFound occurs when path cannot be resolved
type ErrPathNotFound struct {
	wrapErr
	Path string
}

// Error returns a description of the given error as a string
func (e ErrPathNotFound) Error() string {
	return fmt.Sprintf("unable to resolve path, no files found '%s'", e.Path)
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
