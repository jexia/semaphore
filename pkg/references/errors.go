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

func (e ErrUnresolvedFlow) Error() string {
	return fmt.Sprintf("failed to resolve flow '%s'", e.Name)
}

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
	Expression specs.Expression
	// Reference points to the reference which caused the error
	Reference  *specs.PropertyReference
	Breakpoint string
	Path       string
}

func (e ErrUndefinedReference) Error() string {
	return fmt.Sprintf("undefined reference '%s' in '%s.%s'", e.Reference, e.Breakpoint, e.Path)
}

func (e ErrUndefinedReference) Prettify() prettyerr.Error {
	details := map[string]interface{}{
		"Reference":  e.Reference,
		"Breakpoint": e.Breakpoint,
		"Path":       e.Path,
	}

	if e.Expression != nil {
		details["Expression"] = e.Expression.Position()
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

func (e ErrUndefinedResource) Error() string {
	return fmt.Sprintf("undefined resource '%s' in '%s'", e.Reference, e.Breakpoint)
}

func (e ErrUndefinedResource) Prettify() prettyerr.Error {
	var availableRefs []string
	for k, _ := range e.AvailableReferences {
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

type ErrUnresolvedOnError struct {
	wrapErr
	OnError *specs.OnError
}

func (e ErrUnresolvedOnError) Error() string {
	return "failed to resolve OnError"
}

func (e ErrUnresolvedOnError) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UnresolvedOnError",
		Message: e.Error(),
		Details: map[string]interface{}{
			"OnError": e.OnError,
		},
	}
}

type ErrUnresolvedNode struct {
	wrapErr
	Node *specs.Node
}

func (e ErrUnresolvedNode) Error() string {
	return "failed to resolve node"
}

func (e ErrUnresolvedNode) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UnresolvedNode",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Node": e.Node,
		},
	}
}

type ErrUnresolvedParameterMap struct {
	wrapErr
	Parameter *specs.ParameterMap
}

func (e ErrUnresolvedParameterMap) Error() string {
	return "failed to resolve map parameter %s"
}

func (e ErrUnresolvedParameterMap) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UnresolvedParameterMap",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Parameter": e.Parameter,
		},
	}
}

type ErrUnresolvedProperty struct {
	wrapErr
	Property *specs.Property
}

func (e ErrUnresolvedProperty) Error() string {
	return "failed to resolve property"
}

func (e ErrUnresolvedProperty) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UnresolvedProperty",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Property": e.Property,
		},
	}
}

type ErrUnresolvedParams struct {
	wrapErr
	Params map[string]*specs.Property
}

func (e ErrUnresolvedParams) Error() string {
	return "failed to resolve params"
}

func (e ErrUnresolvedParams) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UnresolvedParams",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Params": e.Params,
		},
	}
}

type ErrUnresolvedCall struct {
	wrapErr
	Call *specs.Call
}

func (e ErrUnresolvedCall) Error() string {
	return "failed to resolve call"
}

func (e ErrUnresolvedCall) Prettify() prettyerr.Error {
	return prettyerr.Error{
		Code:    "UnresolvedCall",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Call": e.Call,
		},
	}
}
