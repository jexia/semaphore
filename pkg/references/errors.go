package references

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/lookup"
	"github.com/jexia/semaphore/pkg/specs"
)

type WrapError struct {
	Inner error
}

func (i WrapError) Unwrap() error {
	return i.Inner
}

// ErrUndefinedReference occurs when resolving meets unknown reference.
type ErrUndefinedReference struct {
	WrapError
	Expression specs.Expression
	// Reference points to the reference which caused the error
	Reference  *specs.PropertyReference
	Breakpoint string
	Path       string
}

func (e ErrUndefinedReference) Error() string {
	return fmt.Sprintf("undefined reference '%s' in '%s.%s'", e.Reference, e.Breakpoint, e.Path)
}

// ErrUndefinedResource occurs when resolving meets unknown resource
type ErrUndefinedResource struct {
	WrapError
	// Reference points to the reference which caused the error
	Reference  *specs.PropertyReference
	Breakpoint string
	// AvailableReferences contains the whole list of known references
	AvailableReferences map[string]lookup.ReferenceMap
}

func (e ErrUndefinedResource) Error() string {
	return fmt.Sprintf("undefined resource '%s' in '%s'", e.Reference, e.Breakpoint)
}

type ErrUnresolvedFlow struct {
	WrapError
	Name string
}

func (e ErrUnresolvedFlow) Error() string {
	return fmt.Sprintf("failed to resolve flow '%s'", e.Name)
}

type ErrUnresolvedOnError struct {
	WrapError
	OnError *specs.OnError
}

func (e ErrUnresolvedOnError) Error() string {
	return "failed to resolve OnError"
}

type ErrUnresolvedNode struct {
	WrapError
	Node *specs.Node
}

func (e ErrUnresolvedNode) Error() string {
	return "failed to resovle node"
}

type ErrUnresolvedParameterMap struct {
	WrapError
	Parameter *specs.ParameterMap
}

func (e ErrUnresolvedParameterMap) Error() string {
	return "failed to resolve map parameter"
}

type ErrUnresolvedProperty struct {
	WrapError
	Property *specs.Property
}

func (e ErrUnresolvedProperty) Error() string {
	return "failed to resolve property"
}

type ErrUnresolvedCall struct {
	WrapError
	Call *specs.Call
}

func (e ErrUnresolvedCall) Error() string {
	return "failed to resolve call"
}

type ErrUnresolvedParams struct {
	WrapError
	Params map[string]*specs.Property
}

func (e ErrUnresolvedParams) Error() string {
	return "failed to resolve params"
}
