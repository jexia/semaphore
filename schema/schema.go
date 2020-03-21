package schema

import (
	"github.com/jexia/maestro/specs/types"
)

// Options represents a collection key values
type Options map[string]string

// Collection represents a collection of schemas.
type Collection interface {
	GetService(name string) Service
	GetServices() []Service
	GetMessage(name string) Property
	GetMessages() []Property
}

// Service represents a service which could be called in one of the flows
type Service interface {
	GetComment() string
	GetPackage() string
	GetFullyQualifiedName() string
	GetName() string
	GetHost() string
	GetProtocol() string
	GetCodec() string
	GetMethod(name string) Method
	GetMethods() Methods
	GetOptions() Options
}

// Methods represens a collection of methods
type Methods []Method

// Get attempts to return a method with the given name
func (collection Methods) Get(name string) Method {
	for _, method := range collection {
		if method.GetName() == name {
			return method
		}
	}

	return nil
}

// Method represents a service method
type Method interface {
	GetComment() string
	GetName() string
	GetInput() Property
	GetOutput() Property
	GetOptions() Options
}

// Property represents a object field
type Property interface {
	GetName() string
	GetComment() string
	GetPosition() int32
	GetType() types.Type
	GetLabel() types.Label
	GetNested() map[string]Property
	GetOptions() Options
}

// Resolver when called collects the available schema(s) with the configured configuration
type Resolver func(*Store) error
