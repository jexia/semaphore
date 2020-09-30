package specs

import (
	"errors"
	"fmt"
)

// Repeated represents an array type of fixed size.
type Repeated []Template

// Template returns a Template for a given array. It checks the types of internal
// elements to detect the data type(s).
// Note that all the references must be resolved before calling this method.
func (repeated Repeated) Template() (Template, error) {
	var template Template

	// check if all element types are the same
	// TODO: remove once "oneOf" support is added
	for position := range repeated {
		if position == 0 {
			template = repeated[position]

			continue
		}

		if err := template.Compare(repeated[position]); err != nil {
			return Template{}, fmt.Errorf("all the elements inside the array must have the same type: %w", err)
		}
	}

	// get rid of default value if scalar type
	if template.Scalar != nil {
		template = template.Clone()
		template.Scalar.Default = nil
	}

	return template, nil
}

// Clone repeated.
func (repeated Repeated) Clone() Repeated {
	var clone = make([]Template, len(repeated))

	for index, template := range repeated {
		clone[index] = template.Clone()
	}

	return clone
}

// Compare given repeated to the provided one returning the first mismatch.
func (repeated Repeated) Compare(expected Repeated) error {
	if expected == nil && repeated == nil {
		return nil
	}

	if expected == nil && repeated != nil {
		return errors.New("expected to be nil")
	}

	if expected != nil && repeated == nil {
		return fmt.Errorf("expected to be an array, got %v", nil)
	}

	if len(expected) != len(repeated) {
		return fmt.Errorf("expected to have %d elements, got %d", len(expected), len(repeated))
	}

	left, err := repeated.Template()
	if err != nil {
		return fmt.Errorf("unkown repeated property template: %w", err)
	}

	right, err := expected.Template()
	if err != nil {
		return fmt.Errorf("unkown expected property template: %w", err)
	}

	err = left.Compare(right)
	if err != nil {
		return fmt.Errorf("repeated property: %w", err)
	}

	return nil
}
