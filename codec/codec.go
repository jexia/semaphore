package codec

import (
	"io"

	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/schema"
)

// New constructs a new
type New func(schema.Method) Manager

// Manager represents a marshal/unmarshal codec for a given caller
type Manager interface {
	Marshal(*refs.Store) (io.Reader, error)
	Unmarshal(io.Reader, *refs.Store) error
}
