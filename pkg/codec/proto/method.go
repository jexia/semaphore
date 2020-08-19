package proto

import "github.com/jexia/semaphore/pkg/specs"

// Method represents a service mthod.
type Method interface {
	GetName() string
	GetRequest() map[string]*specs.Property
	GetResponse() map[string]*specs.Property
}

// Methods represents a collection of methods.
type Methods map[string]Method
