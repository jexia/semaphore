package references

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/lookup"
	"github.com/jexia/semaphore/pkg/prettyerr"
	"github.com/jexia/semaphore/pkg/specs"
)

type wrapErr struct {
	Inner error
}

func (i wrapErr) Unwrap() error {
	return i.Inner
}

// ErrUnresolvedFlow occurs when the whole flow cannot be resolved. It's the root error in this package.
type ErrUnresolvedFlow struct {
	wrapErr
	Name string
}

// Error returns a description of the given error as a string
func (e ErrUnresolvedFlow) Error() string {
	return fmt.Sprintf("failed to resolve flow '%s'", e.Name)
}

// Prettify returns the prettified version of the given error
func (e ErrUnresolvedFlow) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UnresolvedFlow",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Name": e.Name,
		},
	}
}

// ErrUndefinedReference occurs when resolving meets unknown reference.
type ErrUndefinedReference struct {
	wrapErr
	Property   *specs.Property
	Breakpoint string
}

// NewErrUndefinedReference constructs a new error thrown when a undefined reference is given.
func NewErrUndefinedReference(inner error, property *specs.Property, breakpoint string) ErrUndefinedReference {
	return ErrUndefinedReference{
		wrapErr: wrapErr{
			Inner: inner,
		},
		Property:   property,
		Breakpoint: breakpoint,
	}
}

// Error returns a description of the given error as a string
func (e ErrUndefinedReference) Error() string {
	return fmt.Sprintf("undefined reference '%s' in '%s.%s'", e.Property.Reference, e.Breakpoint, e.Property.Path)
}

// Prettify returns the prettified version of the given error
func (e ErrUndefinedReference) Prettify() prettyerr.Error {
	details := map[string]interface{}{
		"Reference":  e.Property.Reference,
		"Breakpoint": e.Breakpoint,
		"Path":       e.Property.Path,
	}

	if e.Property.Expr != nil {
		details["Expression"] = e.Property.Expr.Position()
	}

	return prettyerr.Error{
		Code:    "UndefinedReference",
		Message: e.Error(),
		Details: details,
	}
}

// ErrUndefinedResource occurs when resolving meets unknown resource
type ErrUndefinedResource struct {
	wrapErr
	// Reference points to the reference which caused the error
	Reference  *specs.PropertyReference
	Breakpoint string
	// AvailableReferences contains the whole list of known references
	AvailableReferences map[string]lookup.ReferenceMap
}

// Error returns a description of the given error as a string
func (e ErrUndefinedResource) Error() string {
	return fmt.Sprintf("undefined resource '%s' in '%s'", e.Reference, e.Breakpoint)
}

// Prettify returns the prettified version of the given error
func (e ErrUndefinedResource) Prettify() prettyerr.Error {
	var availableRefs []string
	for k := range e.AvailableReferences {
		availableRefs = append(availableRefs, k)
	}

	return prettyerr.Error{
		Code:    "UndefinedResource",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Reference":       e.Reference,
			"Breakpoint":      e.Breakpoint,
			"KnownReferences": availableRefs,
		},
	}
}

// ErrUnresolvedOnError is thrown when Semaphore is unable to resolve a reference
// inside the given on error.
type ErrUnresolvedOnError struct {
	wrapErr
	OnError *specs.OnError
}

// Error returns a description of the given error as a string
func (e ErrUnresolvedOnError) Error() string {
	return "failed to resolve OnError"
}

// Prettify returns the prettified version of the given error
func (e ErrUnresolvedOnError) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UnresolvedOnError",
		Message: e.Error(),
		Details: map[string]interface{}{
			"OnError": e.OnError,
		},
	}
}

// ErrUnresolvedNode is thrown when Semaphore is unable to resolve a reference
// inside the given node.
type ErrUnresolvedNode struct {
	wrapErr
	Node *specs.Node
}

// Error returns a description of the given error as a string
func (e ErrUnresolvedNode) Error() string {
	return "failed to resolve node"
}

// Prettify returns the prettified version of the given error
func (e ErrUnresolvedNode) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UnresolvedNode",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Node": e.Node,
		},
	}
}

// ErrUnresolvedParameterMap is thrown when Semaphore is unable to resolve
// a reference in the given parameter map.
type ErrUnresolvedParameterMap struct {
	wrapErr
	Parameter *specs.ParameterMap
}

// Error returns a description of the given error as a string
func (e ErrUnresolvedParameterMap) Error() string {
	return "failed to resolve map parameter"
}

// Prettify returns the prettified version of the given error
func (e ErrUnresolvedParameterMap) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UnresolvedParameterMap",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Parameter": e.Parameter,
		},
	}
}

// ErrUnresolvedProperty is thrown when Semaphore is unable to resolve a reference
// inside the given property.
type ErrUnresolvedProperty struct {
	wrapErr
	Property *specs.Property
}

// NewErrUnresolvedProperty is thrown when the given property has not been resolved
func NewErrUnresolvedProperty(inner error, property *specs.Property) ErrUnresolvedProperty {
	return ErrUnresolvedProperty{
		wrapErr: wrapErr{
			Inner: inner,
		},
		Property: property,
	}
}

func (e ErrUnresolvedProperty) Error() string {
	return "failed to resolve property"
}

// Prettify returns the prettified version of the given error
func (e ErrUnresolvedProperty) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UnresolvedProperty",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Property": e.Property,
		},
	}
}

// ErrUnresolvedParams is thrown when Semaphore is unable to resolve a reference
// inside the given parameters.
type ErrUnresolvedParams struct {
	wrapErr
	Params map[string]*specs.Property
}

func (e ErrUnresolvedParams) Error() string {
	return "failed to resolve params"
}

// Prettify returns the prettified version of the given error
func (e ErrUnresolvedParams) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UnresolvedParams",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Params": e.Params,
		},
	}
}

// ErrUnresolvedCall is thrown when Semaphore is unable to resolve a reference
// inside the given call.
type ErrUnresolvedCall struct {
	wrapErr
	Call *specs.Call
}

func (e ErrUnresolvedCall) Error() string {
	return "failed to resolve call"
}

// Prettify returns the prettified version of the given error
func (e ErrUnresolvedCall) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UnresolvedCall",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Call": e.Call,
		},
	}
}
