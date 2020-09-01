package references

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/lookup"
	"github.com/jexia/semaphore/pkg/pretty_errors"
	"github.com/jexia/semaphore/pkg/specs"
)

type wrapErr struct {
	Inner error
}

func (i wrapErr) Unwrap() error {
	return i.Inner
}

type (
	// ErrUnresolvedFlow occurs when the whole flow cannot be resolved. It's the root error in this package.
	ErrUnresolvedFlow struct {
		wrapErr
		Name string
	}

	// ErrUndefinedReference occurs when resolving meets unknown reference.
	ErrUndefinedReference struct {
		wrapErr
		Expression specs.Expression
		// Reference points to the reference which caused the error
		Reference  *specs.PropertyReference
		Breakpoint string
		Path       string
	}

	// ErrUndefinedResource occurs when resolving meets unknown resource
	ErrUndefinedResource struct {
		wrapErr
		// Reference points to the reference which caused the error
		Reference  *specs.PropertyReference
		Breakpoint string
		// AvailableReferences contains the whole list of known references
		AvailableReferences map[string]lookup.ReferenceMap
	}

	ErrUnresolvedOnError struct {
		wrapErr
		OnError *specs.OnError
	}

	ErrUnresolvedNode struct {
		wrapErr
		Node *specs.Node
	}

	ErrUnresolvedParameterMap struct {
		wrapErr
		Parameter *specs.ParameterMap
	}

	ErrUnresolvedProperty struct {
		wrapErr
		Property *specs.Property
	}

	ErrUnresolvedParams struct {
		wrapErr
		Params map[string]*specs.Property
	}

	ErrUnresolvedCall struct {
		wrapErr
		Call *specs.Call
	}
)

func (e ErrUndefinedReference) Error() string {
	return fmt.Sprintf("undefined reference '%s' in '%s.%s'", e.Reference, e.Breakpoint, e.Path)
}
func (e ErrUndefinedReference) Prettify() pretty_errors.Error {
	details := map[string]interface{}{
		"Reference":  e.Reference,
		"Breakpoint": e.Breakpoint,
		"Path":       e.Path,
	}

	if e.Expression != nil {
		details["Expression"] = e.Expression.Position()
	}

	return pretty_errors.Error{
		Code:    "UndefinedReference",
		Message: e.Error(),
		Details: details,
	}
}

func (e ErrUndefinedResource) Error() string {
	return fmt.Sprintf("undefined resource '%s' in '%s'", e.Reference, e.Breakpoint)
}
func (e ErrUndefinedResource) Prettify() pretty_errors.Error {
	var availableRefs []string
	for k, _ := range e.AvailableReferences {
		availableRefs = append(availableRefs, k)
	}

	return pretty_errors.Error{
		Code:    "UndefinedResource",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Reference":  e.Reference,
			"Breakpoint": e.Breakpoint,
			"KnownReferences": availableRefs,
		},
	}
}

func (e ErrUnresolvedFlow) Error() string {
	return fmt.Sprintf("failed to resolve flow '%s'", e.Name)
}
func (e ErrUnresolvedFlow) Prettify() pretty_errors.Error {
	return pretty_errors.Error{
		Code:    "UnresolvedFlow",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Name": e.Name,
		},
	}
}

func (e ErrUnresolvedOnError) Error() string {
	return "failed to resolve OnError"
}
func (e ErrUnresolvedOnError) Prettify() pretty_errors.Error {
	return pretty_errors.Error{
		Code:    "UnresolvedOnError",
		Message: e.Error(),
		Details: map[string]interface{}{
			"OnError": e.OnError,
		},
	}
}

func (e ErrUnresolvedNode) Error() string {
	return "failed to resovle node"
}
func (e ErrUnresolvedNode) Prettify() pretty_errors.Error {
	return pretty_errors.Error{
		Code:    "UnresolvedNode",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Node": e.Node,
		},
	}
}

func (e ErrUnresolvedParameterMap) Error() string {
	return "failed to resolve map parameter %s"
}
func (e ErrUnresolvedParameterMap) Prettify() pretty_errors.Error {
	return pretty_errors.Error{
		Code:    "UnresolvedParameterMap",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Parameter": e.Parameter,
		},
	}
}

func (e ErrUnresolvedProperty) Error() string {
	return "failed to resolve property"
}
func (e ErrUnresolvedProperty) Prettify() pretty_errors.Error {
	return pretty_errors.Error{
		Code:    "UnresolvedProperty",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Property": e.Property,
		},
	}
}

func (e ErrUnresolvedCall) Error() string {
	return "failed to resolve call"
}
func (e ErrUnresolvedCall) Prettify() pretty_errors.Error {
	return pretty_errors.Error{
		Code:    "UnresolvedCall",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Call": e.Call,
		},
	}
}

func (e ErrUnresolvedParams) Error() string {
	return "failed to resolve params"
}
func (e ErrUnresolvedParams) Prettify() pretty_errors.Error {
	return pretty_errors.Error{
		Code:    "UnresolvedParams",
		Message: e.Error(),
		Details: map[string]interface{}{
			"Params": e.Params,
		},
	}
}
