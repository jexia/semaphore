package specs

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jexia/semaphore/pkg/specs/types"
)

// Scalar value.
type Scalar struct {
	Default interface{} `json:"default,omitempty"`
	Type    types.Type  `json:"type,omitempty"`
}

// UnmarshalJSON corrects the 64bit data types in accordance with golang
func (scalar *Scalar) UnmarshalJSON(data []byte) error {
	if scalar == nil {
		return nil
	}

	type sc Scalar
	t := sc{}

	err := json.Unmarshal(data, &t)
	if err != nil {
		return err
	}

	*scalar = Scalar(t)
	scalar.Clean()

	return nil
}

// Clean fixes the type casting issue of unmarshal
func (scalar *Scalar) Clean() {
	if scalar.Default == nil {
		return
	}

	switch scalar.Type {
	case types.Int64, types.Sint64, types.Sfixed64:
		_, ok := scalar.Default.(int64)
		if !ok {
			scalar.Default = int64(scalar.Default.(float64))
		}
	case types.Uint64, types.Fixed64:
		_, ok := scalar.Default.(uint64)
		if !ok {
			scalar.Default = uint64(scalar.Default.(float64))
		}
	case types.Int32, types.Sint32, types.Sfixed32:
		_, ok := scalar.Default.(int32)
		if !ok {
			scalar.Default = int32(scalar.Default.(float64))
		}
	case types.Uint32, types.Fixed32:
		_, ok := scalar.Default.(uint32)
		if !ok {
			scalar.Default = uint32(scalar.Default.(float64))
		}
	}
}

// Clone scalar value.
func (scalar Scalar) Clone() *Scalar {
	return &Scalar{
		Default: scalar.Default,
		Type:    scalar.Type,
	}
}

// Compare the given scalar against the expected and return the first met difference
// as an error.
func (scalar *Scalar) Compare(expected *Scalar) error {
	if expected == nil && scalar == nil {
		return nil
	}

	if expected == nil && scalar != nil {
		return errors.New("expected to be nil")
	}

	if expected != nil && scalar == nil {
		return fmt.Errorf("expected to be %v, got %v", expected.Type, nil)
	}

	if expected.Type != scalar.Type {
		return fmt.Errorf("expected to be %v, got %v", expected.Type, scalar.Type)
	}

	return nil
}
