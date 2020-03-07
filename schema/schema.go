package schema

import "github.com/jexia/maestro/specs/types"

// Options represents a collection key values
type Options map[string]string

// Collection represents a collection of schemas.
type Collection interface {
	GetService(name string) Service
	GetProperty(name string) Property
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
	GetInput() Property
	GetOutput() Property
	GetOptions() Options
}

// Property represents a object field
type Property interface {
	GetName() string
	GetType() types.Type
	GetLabel() types.Label
	GetNested() map[string]Property
	GetOptions() Options
}
