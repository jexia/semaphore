package graphql

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/prettyerr"
	"github.com/jexia/semaphore/pkg/specs/types"
)

type wrapErr struct {
	Inner error
}

func (i wrapErr) Unwrap() error {
	return i.Inner
}

// ErrUnexpectedType occurs when a non-object is passed as property
type ErrUnexpectedType struct {
	wrapErr
	Type     types.Type
	Expected types.Type
}

// Error returns a description of the given error as a string
func (e ErrUnexpectedType) Error() string {
	return fmt.Sprintf("unexpected property type, received '%s', expected '%s'", e.Type, e.Expected)
}

// Prettify returns the prettified version of the given error
func (e ErrUnexpectedType) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UnexpectedType",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Type":     e.Type,
			"Expected": e.Expected,
		},
	}
}

// ErrTypeMismatch occurs when there is a mismatch in schema definations
type ErrTypeMismatch struct {
	wrapErr
	Expected string
	Type     string
}

// Error returns a description of the given error as a string
func (e ErrTypeMismatch) Error() string {
	return fmt.Sprintf("unable set field '%s' in '%s'", e.Type, e.Expected)
}

// Prettify returns the prettified version of the given error
func (e ErrTypeMismatch) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "TypeMismatch",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Type":     e.Type,
			"Expected": e.Expected,
		},
	}
}

// ErrFieldAlreadySet occurs when field has already being set
type ErrFieldAlreadySet struct {
	wrapErr
	Field string
	Path  string
}

// Error returns a description of the given error as a string
func (e ErrFieldAlreadySet) Error() string {
	return fmt.Sprintf("field already set '%s' in '%s'", e.Field, e.Path)
}

// Prettify returns the prettified version of the given error
func (e ErrFieldAlreadySet) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "FieldAlreadySet",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Field": e.Field,
			"Path":  e.Path,
		},
	}
}

// ErrDuplicateObject occurs when a duplicate object is provided
type ErrDuplicateObject struct {
	wrapErr
	Name string
}

// Error returns a description of the given error as a string
func (e ErrDuplicateObject) Error() string {
	return fmt.Sprintf("duplicate object '%s'", e.Name)
}

// Prettify returns the prettified version of the given error
func (e ErrDuplicateObject) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "DuplicateObject",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Name": e.Name,
		},
	}
}

// ErrUnknownBase occurs when an unknown base is provided
type ErrUnknownBase struct {
	wrapErr
	Base string
}

// Error returns a description of the given error as a string
func (e ErrUnknownBase) Error() string {
	return fmt.Sprintf("unknown base '%s', expected query or mutation", e.Base)
}

// Prettify returns the prettified version of the given error
func (e ErrUnknownBase) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UnknownBase",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Base": e.Base,
		},
	}
}
