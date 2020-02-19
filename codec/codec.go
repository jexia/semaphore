package codec

import (
	"io"

	"github.com/jexia/maestro/refs"
)

// Manager represents a marshal/unmarshal codec for a given caller
type Manager interface {
	Marshal(*refs.Store) (io.Reader, error)
	Unmarshal(io.Reader, *refs.Store) error
}
