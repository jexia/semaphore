package schema

import "github.com/jexia/maestro/specs/types"

// Options represents a collection key values
type Options map[string]string

// Collection represents a collection of schemas.
type Collection interface {
	GetService(name string) Service
}

// Service represents a service which could be called in one of the flows
type Service interface {
	GetName() string
	GetMethod(name string) Method
	GetMethods() []Method
	GetOptions() Options
}

// Method represents a service method
type Method interface {
	GetName() string
	GetInput() Object
	GetOutput() Object
	GetOptions() Options
}

// Object represents a data object
type Object interface {
	GetFields() []Field
	GetField(name string) Field
	GetOptions() Options
}

// Field represents a object field
type Field interface {
	GetName() string
	GetType() types.Type
	GetLabel() types.Label
	GetObject() Object
	GetOptions() Options
}
