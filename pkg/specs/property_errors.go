package specs

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/prettyerr"
	"github.com/jexia/semaphore/pkg/specs/labels"
)

// ErrUndefinedSchema occurs when schema is not found at path
type ErrUndefinedSchema struct {
	Path string
	Expr Expression
}

// Error returns a description of the given error as a string
func (e ErrUndefinedSchema) Error() string {
	message := fmt.Sprintf("unable to check types for '%s' no schema given", e.Path)
	if e.Expr == nil {
		return message
	}
	return fmt.Sprintf("%s %s", e.Expr.Position(), message)
}

// Prettify returns the prettified version of the given error
func (e ErrUndefinedSchema) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UndefinedSchema",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Path": e.Path,
		},
	}
}

// ErrTypeMismatch occurs when given typs does not match with expected type
type ErrTypeMismatch struct {
	Type     interface{}
	Expected interface{}
	Path     string
	Expr     Expression
}

// Error returns a description of the given error as a string
func (e ErrTypeMismatch) Error() string {
	message := fmt.Sprintf("cannot use type (%s) for '%s', expected (%s)", e.Type, e.Path, e.Expected)
	if e.Expr == nil {
		return message
	}
	return fmt.Sprintf("%s %s", e.Expr.Position(), message)
}

// Prettify returns the prettified version of the given error
func (e ErrTypeMismatch) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "TypeMismatch",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Path":     e.Path,
			"Type":     e.Type,
			"Expected": e.Expected,
		},
	}
}

// ErrLabelMismatch occurs when given label does not match with expected label
type ErrLabelMismatch struct {
	Label    labels.Label
	Expected labels.Label
	Path     string
	Expr     Expression
}

// Error returns a description of the given error as a string
func (e ErrLabelMismatch) Error() string {
	message := fmt.Sprintf("cannot use label (%s) for '%s', expected (%s)", e.Label, e.Path, e.Expected)
	if e.Expr == nil {
		return message
	}
	return fmt.Sprintf("%s %s", e.Expr.Position(), message)
}

// Prettify returns the prettified version of the given error
func (e ErrLabelMismatch) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "LabelMismatch",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Path":     e.Path,
			"Label":    e.Label,
			"Expected": e.Expected,
		},
	}
}

// ErrUndeclaredSchema occurs when nested object does not have schema
type ErrUndeclaredSchema struct {
	Name string
	Path string
	Expr Expression
}

// Error returns a description of the given error as a string
func (e ErrUndeclaredSchema) Error() string {
	message := fmt.Sprintf("property '%s' has a nested object but schema does not '%s'", e.Path, e.Name)
	if e.Expr == nil {
		return message
	}
	return fmt.Sprintf("%s %s", e.Expr.Position(), message)
}

// Prettify returns the prettified version of the given error
func (e ErrUndeclaredSchema) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UndeclaredSchema",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Path": e.Path,
			"Name": e.Name,
		},
	}
}

// ErrUndeclaredSchemaInProperty occurs when nested property does not have schema
type ErrUndeclaredSchemaInProperty struct {
	Name string
	Path string
	Expr Expression
}

// Error returns a description of the given error as a string
func (e ErrUndeclaredSchemaInProperty) Error() string {
	message := fmt.Sprintf("undefined schema nested message property '%s' in flow '%s'", e.Path, e.Name)
	if e.Expr == nil {
		return message
	}
	return fmt.Sprintf("%s %s", e.Expr.Position(), message)
}

// Prettify returns the prettified version of the given error
func (e ErrUndeclaredSchemaInProperty) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UndeclaredSchemaInProperty",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Path": e.Path,
			"Name": e.Name,
		},
	}
}
