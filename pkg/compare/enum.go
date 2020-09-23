package compare

import (
	"errors"
	"fmt"

	"github.com/jexia/semaphore/pkg/specs"
)

// CompareEnumValues compares the given enum value against of the expected one and returns the first met difference as error.
func CompareEnumValues(expected *specs.EnumValue, given *specs.EnumValue) error {
	if expected == nil && given == nil {
		return nil
	}

	if expected == nil && given != nil {
		return errors.New("expected to be nil")
	}

	if expected != nil && given == nil {
		return fmt.Errorf("expected to be %v:%v, got %v", expected.Key, expected.Position, nil)
	}

	if expected.Key != given.Key || expected.Position != given.Position {
		return fmt.Errorf("expected to be %v:%v, got %v:%v", expected.Key, expected.Position, given.Key, given.Position)
	}

	return nil
}

// CompareEnums compares the given enum against of the expected one and returns the first met difference as error.
func CompareEnums(expected *specs.Enum, given *specs.Enum) error {
	if expected == nil && given == nil {
		return nil
	}

	if expected == nil && given != nil {
		return errors.New("expected to be nil")
	}

	if expected != nil && given == nil {
		return fmt.Errorf("expected to be %v, got %v", expected.Name, nil)
	}

	if expected.Name != given.Name {
		return fmt.Errorf("expected to be %v, got %v", expected.Name, given.Name)
	}

	if len(expected.Keys) != len(given.Keys) {
		return fmt.Errorf("expected to have %v keys, got %v", len(expected.Keys), given.Keys)
	}

	if len(expected.Positions) != len(given.Positions) {
		return fmt.Errorf("expected to have %v positions, got %v", len(expected.Positions), given.Positions)
	}

	for expectedKey, expectedValue := range expected.Keys {
		// given enum does not include the current key
		givenEnumValue, ok := given.Keys[expectedKey]
		if !ok {
			return fmt.Errorf("expected to have %v key", expectedKey)
		}

		err := CompareEnumValues(expectedValue, givenEnumValue)
		if err != nil {
			return fmt.Errorf("value mismatch: %w", err)
		}
	}

	for expectedPos, expectedValue := range expected.Positions {
		// given enum does not include the current position
		givenEnumValue, ok := given.Positions[expectedPos]
		if !ok {
			return fmt.Errorf("expected to have %v position", expectedPos)
		}

		err := CompareEnumValues(expectedValue, givenEnumValue)
		if err != nil {
			return fmt.Errorf("value mismatch: %w", err)
		}
	}

	return nil
}
