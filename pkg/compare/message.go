package compare

import (
	"errors"
	"fmt"

	"github.com/jexia/semaphore/pkg/specs"
)

// CompareMessages compares the given message against of the expected one and returns the first met difference as error.
func CompareMessages(expected specs.Message, given specs.Message) error {
	if expected == nil && given == nil {
		return nil
	}

	if expected == nil && given != nil {
		return errors.New("expected to be nil")
	}

	if expected != nil && given == nil {
		return fmt.Errorf("expected to be an object, got %v", nil)
	}

	if len(expected) != len(given) {
		return fmt.Errorf("expected to have %v properties, got %v", len(expected), len(given))
	}

	for expectedName, expectedProperty := range expected {
		// given message does not include the current property
		givenProperty, ok := given[expectedName]
		if !ok {
			return fmt.Errorf("expected the object has field '%s'", expectedName)
		}

		err := CheckPropertyTypes(expectedProperty, givenProperty)
		if err != nil {
			return fmt.Errorf("object property mismatch: %w", err)
		}
	}

	return nil
}
