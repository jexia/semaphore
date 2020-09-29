package specs

import (
	"errors"
	"fmt"

	"github.com/jexia/semaphore/pkg/specs/metadata"
)

// Enum represents a enum configuration
type Enum struct {
	*metadata.Meta
	Name        string                `json:"name,omitempty" yaml:"name,omitempty"`
	Description string                `json:"description,omitempty" yaml:"description,omitempty"`
	Keys        map[string]*EnumValue `json:"keys,omitempty" yaml:"keys,omitempty"`
	Positions   map[int32]*EnumValue  `json:"positions,omitempty" yaml:"positions,omitempty"`
}

// Clone enum schema.
func (enum Enum) Clone() *Enum { return &enum }

// Compare the given enum value against of the expected one and return the first
// met difference as error.
func (enum *Enum) Compare(expected *Enum) error {
	if expected == nil && enum == nil {
		return nil
	}

	if expected == nil && enum != nil {
		return errors.New("expected to be nil")
	}

	if expected != nil && enum == nil {
		return fmt.Errorf("expected to be %v, got %v", expected.Name, nil)
	}

	if expected.Name != enum.Name {
		return fmt.Errorf("expected to be %v, got %v", expected.Name, enum.Name)
	}

	if len(expected.Keys) != len(enum.Keys) {
		return fmt.Errorf("expected to have %v keys, got %v", len(expected.Keys), enum.Keys)
	}

	if len(expected.Positions) != len(enum.Positions) {
		return fmt.Errorf("expected to have %v positions, got %v", len(expected.Positions), enum.Positions)
	}

	for expectedKey, expectedValue := range expected.Keys {
		// given enum does not include the current key
		enumValue, ok := enum.Keys[expectedKey]
		if !ok {
			return fmt.Errorf("expected to have %v key", expectedKey)
		}

		err := enumValue.Compare(expectedValue)
		if err != nil {
			return fmt.Errorf("value mismatch: %w", err)
		}
	}

	for expectedPos, expectedValue := range expected.Positions {
		// given enum does not include the current position
		enumValue, ok := enum.Positions[expectedPos]
		if !ok {
			return fmt.Errorf("expected to have %v position", expectedPos)
		}

		err := enumValue.Compare(expectedValue)
		if err != nil {
			return fmt.Errorf("value mismatch: %w", err)
		}
	}

	return nil
}

// EnumValue represents a enum configuration
type EnumValue struct {
	*metadata.Meta
	Key         string `json:"key,omitempty"`
	Position    int32  `json:"position,omitempty"`
	Description string `json:"description,omitempty"`
}

// Compare the given enum value against the expected one and returns the first
// met difference as error.
func (value *EnumValue) Compare(expected *EnumValue) error {
	if expected == nil && value == nil {
		return nil
	}

	if expected == nil && value != nil {
		return errors.New("expected to be nil")
	}

	if expected != nil && value == nil {
		return fmt.Errorf("expected to be %v:%v, got %v", expected.Key, expected.Position, nil)
	}

	if expected.Key != value.Key || expected.Position != value.Position {
		return fmt.Errorf("expected to be %v:%v, got %v:%v", expected.Key, expected.Position, value.Key, value.Position)
	}

	return nil
}
