package flow

import (
	"io"

	"github.com/jexia/maestro/refs"
)

// Codec represents a marshal/unmarshal codec for a given caller
type Codec interface {
	Marshal(*refs.Store) (io.Reader, error)
	Unmarshal(io.Reader, *refs.Store) error
}
