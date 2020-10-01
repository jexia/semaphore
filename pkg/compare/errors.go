package compare

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/prettyerr"
)

// ErrHeaderTypeMismatch occurs when expected header is not revieved
type ErrHeaderTypeMismatch struct {
	Type     interface{}
	Expected interface{}
	Path     string
	Flow     string
}

// Error returns a description of the given error as a string
func (e ErrHeaderTypeMismatch) Error() string {
	return fmt.Sprintf("cannot use type (%s) for 'header.%s' in flow '%s', expected (%s)", e.Type, e.Path, e.Flow, e.Expected)
}

// Prettify returns the prettified version of the given error
func (e ErrHeaderTypeMismatch) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "HeaderTypeMismatch",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Path":     e.Path,
			"Flow":     e.Flow,
			"Type":     e.Type,
			"Expected": e.Expected,
		},
	}
}

// ErrUndefinedObject occurs when flow output object is not defined
type ErrUndefinedObject struct {
	Flow   string
	Schema string
}

// Error returns a description of the given error as a string
func (e ErrUndefinedObject) Error() string {
	return fmt.Sprintf("undefined flow output object '%s' in '%s'", e.Schema, e.Flow)
}

// Prettify returns the prettified version of the given error
func (e ErrUndefinedObject) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UndefinedObject",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Flow":   e.Flow,
			"Schema": e.Schema,
		},
	}
}

// ErrUndefinedService occurs when called service is not defined in flow
type ErrUndefinedService struct {
	Flow    string
	Service string
}

// Error returns a description of the given error as a string
func (e ErrUndefinedService) Error() string {
	return fmt.Sprintf("undefined service '%s' in flow '%s'", e.Service, e.Flow)
}

// Prettify returns the prettified version of the given error
func (e ErrUndefinedService) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UndefinedService",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Flow":    e.Flow,
			"Service": e.Service,
		},
	}
}

// ErrUndefinedMethod occurs when called method is not defined in flow
type ErrUndefinedMethod struct {
	Flow   string
	Method string
}

// Error returns a description of the given error as a string
func (e ErrUndefinedMethod) Error() string {
	return fmt.Sprintf("undefined method '%s' in flow '%s'", e.Method, e.Flow)
}

// Prettify returns the prettified version of the given error
func (e ErrUndefinedMethod) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UndefinedMethod",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Flow":   e.Flow,
			"Method": e.Method,
		},
	}
}
