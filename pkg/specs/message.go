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
	return message.clone(make(map[string]*Template))
}

func (message Message) clone(seen map[string]*Template) Message {
	var clone = make(map[string]*Property, len(message))

	for key := range message {
		clone[key] = message[key].clone(seen)
	}

	return clone
}

// Compare given message to the provided one returning the first mismatch.
func (message Message) Compare(expected Message) error {
	return message.compare(make(map[string]*Template), expected)
}

func (message Message) compare(seen map[string]*Template, expected Message) error {
	if expected == nil && message == nil {
		return nil
	}

	if expected == nil && message != nil {
		return errors.New("expected to be nil")
	}

	if expected != nil && message == nil {
		return fmt.Errorf("expected to be an object, got %v", nil)
	}

	for key, property := range message {
		nested, ok := expected[key]
		if !ok {
			return fmt.Errorf("object has unknown field '%s'", key)
		}

		if err := property.compare(seen, nested); err != nil {
			return fmt.Errorf("object property mismatch: %w", err)
		}
	}

	return nil
}
