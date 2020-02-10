package schema

import "github.com/jexia/maestro/specs"

// Collection represents a collection of schemas.
type Collection interface {
	GetService(name string) Service
}

// Service represents a service which could be called in one of the flows
type Service interface {
	GetName() string
	GetMethod(name string) Method
}

// Method represents a service method
type Method interface {
	GetName() string
	GetInput() Object
	GetOutput() Object
}

// Object represents a data object
type Object interface {
	GetField(name string) Field
}

// Field represents a object field
type Field interface {
	GetName() string
	GetType() specs.Type
	GetLabel() specs.Label
	GetObject() Object
}
