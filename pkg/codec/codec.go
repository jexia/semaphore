package codec

import (
	"io"

	"github.com/jexia/semaphore/pkg/refs"
	"github.com/jexia/semaphore/pkg/specs"
)

// Constructors represent a collection of codec constructors
type Constructors map[string]Constructor

// Get attempts to fetch a codec constructor from the collection matching the given key
func (collection Constructors) Get(key string) Constructor {
	return collection[key]
}

// Constructor is capable of constructing new codec managers for the given resource and specs
type Constructor interface {
	Name() string
	New(resource string, specs *specs.ParameterMap) (Manager, error)
}

// Manager represents a marshal/unmarshal codec for a given caller
type Manager interface {
	Name() string
	Property() *specs.Property
	Marshal(refs.Store) (io.Reader, error)
	Unmarshal(io.Reader, refs.Store) error
}
