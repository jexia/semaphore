package http

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/prettyerr"
)

// ErrRouteConflict is returned when HTTP route conflict is detected.
type ErrRouteConflict string

func (e ErrRouteConflict) Error() string { return string(e) }

type wrapErr struct {
	Inner error
}

func (i wrapErr) Unwrap() error {
	return i.Inner
}

// ErrUnknownMethod occurs when undefined method is called
type ErrUnknownMethod struct {
	Method  string
	Service string
}

// Error returns a description of the given error as a string
func (e ErrUnknownMethod) Error() string {
	return fmt.Sprintf("unknown method '%s' for service '%s'", e.Method, e.Service)
}

// Prettify returns the prettified version of the given error
func (e ErrUnknownMethod) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UnknownMethod",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Method":  e.Method,
			"Service": e.Service,
		},
	}
}

// ErrUndefinedCodec occurs when undefined codec is called
type ErrUndefinedCodec struct {
	Codec string
}

// Error returns a description of the given error as a string
func (e ErrUndefinedCodec) Error() string {
	return fmt.Sprintf("request codec not found '%s'", e.Codec)
}

// Prettify returns the prettified version of the given error
func (e ErrUndefinedCodec) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UndefinedCodec",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Codec": e.Codec,
		},
	}
}

// ErrInvalidHost occurs when provided host is invalid
type ErrInvalidHost struct {
	wrapErr
	Host string
}

// Error returns a description of the given error as a string
func (e ErrInvalidHost) Error() string {
	return fmt.Sprintf("unable to parse the proxy forward host '%s'", e.Host)
}

// Prettify returns the prettified version of the given error
func (e ErrInvalidHost) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "InvalidHost",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Codec": e.Host,
		},
	}
}
