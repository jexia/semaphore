package codec

import (
	"io"

	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs"
)

// New constructs a new message codec handler
type New func(resource string, schema schema.Object, specs specs.Object) (Manager, error)

// Manager represents a marshal/unmarshal codec for a given caller
type Manager interface {
	Marshal(*refs.Store) (io.Reader, error)
	Unmarshal(io.Reader, *refs.Store) error
}
