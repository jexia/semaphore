package references

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/lookup"
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

func (e ErrUndefinedResource) Error() string {
	return fmt.Sprintf("undefined resource '%s' in '%s'", e.Reference, e.Breakpoint)
}

func (e ErrUnresolvedFlow) Error() string {
	return fmt.Sprintf("failed to resolve flow '%s'", e.Name)
}

func (e ErrUnresolvedOnError) Error() string {
	return "failed to resolve OnError"
}

func (e ErrUnresolvedNode) Error() string {
	return "failed to resovle node"
}

func (e ErrUnresolvedParameterMap) Error() string {
	return "failed to resolve map parameter %s"
}

func (e ErrUnresolvedProperty) Error() string {
	return "failed to resolve property"
}

func (e ErrUnresolvedCall) Error() string {
	return "failed to resolve call"
}

func (e ErrUnresolvedParams) Error() string {
	return "failed to resolve params"
}
