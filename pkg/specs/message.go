package specs

import (
	"errors"
	"fmt"
	"sort"
)

// Message represents an object which keeps the original order of keys.
type Message map[string]*Property

// SortedProperties returns the available properties as a properties list
// ordered base on the properties position.
func (message Message) SortedProperties() PropertyList {
	result := make(PropertyList, 0, len(message))

	for _, property := range message {
		result = append(result, property)
	}

	sort.Sort(result)
	return result
}

// Clone the message.
func (message Message) Clone() Message {
	var clone = make(map[string]*Property, len(message))

	for key := range message {
		clone[key] = message[key].Clone()
	}

	return clone
}

// Compare given message to the provided one returning the first mismatch.
func (message Message) Compare(expected Message) error {
	if expected == nil && message == nil {
		return nil
	}

	if expected == nil && message != nil {
		return errors.New("expected to be nil")
	}

	if expected != nil && message == nil {
		return fmt.Errorf("expected to be an object, got %v", nil)
	}

	if len(expected) != len(message) {
		return fmt.Errorf("expected to have %d properties, got %d", len(expected), len(message))
	}

	for expectedName, expectedProperty := range expected {
		// given message does not include the current property
		givenProperty, ok := message[expectedName]
		if !ok {
			return fmt.Errorf("expected the object has field '%s'", expectedName)
		}

		if err := givenProperty.Compare(expectedProperty); err != nil {
			return fmt.Errorf("object property mismatch: %w", err)
		}
	}

	return nil
}
