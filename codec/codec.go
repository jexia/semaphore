package codec

import (
	"io"

	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/specs"
)

// Constructor is capable of constructing new codec managers for the given resource and specs
type Constructor interface {
	Name() string
	New(resource string, specs *specs.ParameterMap) (Manager, error)
}

// Manager represents a marshal/unmarshal codec for a given caller
type Manager interface {
	Marshal(*refs.Store) (io.Reader, error)
	Unmarshal(io.Reader, *refs.Store) error
}
