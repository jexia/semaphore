package providers

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/prettyerr"
)

// ErrUndefinedObject occurs when Object is not found in schema
type ErrUndefinedObject struct {
	Schema string
}

// Error returns a description of the given error as a string
func (e ErrUndefinedObject) Error() string {
	return fmt.Sprintf("object '%s', is unavailable inside the schema collection", e.Schema)
}

// Prettify returns the prettified version of the given error
func (e ErrUndefinedObject) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UnavailableObject",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Schema": e.Schema,
		},
	}
}

// ErrUndefinedService occurs when Service is not defined in a flow
type ErrUndefinedService struct {
	Service string
	Flow    string
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
			"Service": e.Service,
			"Flow":    e.Flow,
		},
	}
}

// ErrUndefinedMethod occurs when method is not defined in a flow
type ErrUndefinedMethod struct {
	Method string
	Flow   string
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
			"Method": e.Method,
			"Flow":   e.Flow,
		},
	}
}

// ErrUndefinedOutput occurs when output is not defined in a flow
type ErrUndefinedOutput struct {
	Output string
	Flow   string
}

// Error returns a description of the given error as a string
func (e ErrUndefinedOutput) Error() string {
	return fmt.Sprintf("undefined method output property '%s' in flow '%s'", e.Output, e.Flow)
}

// Prettify returns the prettified version of the given error
func (e ErrUndefinedOutput) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UndefinedOutput",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Output": e.Output,
			"Flow":   e.Flow,
		},
	}
}

// ErrUndefinedProperty occurs when property is not defined in a flow
type ErrUndefinedProperty struct {
	Property string
	Flow     string
}

// Error returns a description of the given error as a string
func (e ErrUndefinedProperty) Error() string {
	return fmt.Sprintf("undefined schema nested message property '%s' in flow '%s'", e.Property, e.Flow)
}

// Prettify returns the prettified version of the given error
func (e ErrUndefinedProperty) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UndefinedProperty",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Property": e.Property,
			"Flow":     e.Flow,
		},
	}
}
