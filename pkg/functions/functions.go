package functions

import (
	"github.com/jexia/maestro/pkg/refs"
	"github.com/jexia/maestro/pkg/specs"
)

// Custom represents a collection of custom defined functions that could be called inside a template
type Custom map[string]Intermediate

// Intermediate prepares the custom defined function.
// The given arguments represent the exprected types that are passed when called.
// Properties returned should be absolute.
type Intermediate func(args ...*specs.Property) (*specs.Property, Exec, error)

// Exec is a executable function.
// A store should be returned which could be used to encode the function property
type Exec func(store refs.Store) error

// Collection represents a collection of stacks grouped by nodes
type Collection map[*specs.Node]Stack

// Reserve reserves a new function stack for the given node.
// If a stack already exists for the given node is it returned.
func (collection Collection) Reserve(node *specs.Node) Stack {
	stack, has := collection[node]
	if has {
		return stack
	}

	collection[node] = Stack{}
	return collection[node]
}

// Stack represents a collection of functions
type Stack map[string]*Function

// Function represents a custom defined function
type Function struct {
	Arguments []*specs.Property
	Fn        Exec
	Returns   *specs.Property
}
