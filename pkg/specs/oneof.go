package specs

import (
	"errors"
	"fmt"
)

// OneOf is a mixed type to let the schema validate values against exactly one of the templates.
// Example:
// OneOf{
//   {Scalar: &Scalar{Type: types.String}},
//   {Scalar: &Scalar{Type: types.Int32}},
//   {Message: &Message{...}},
// }
// A given value must be one of these types: string, int32 or the message.
type OneOf []*Property

func (oneOf OneOf) String() string { return dump(oneOf) }

// Clone OneOf.
func (oneOf OneOf) Clone() OneOf {
	var clone = make(OneOf, len(oneOf))

	for key := range oneOf {
		clone[key] = oneOf[key].Clone()
	}

	return clone
}

// Compare checks whether given OneOf mathches the provided one.
func (oneOf OneOf) Compare(expected OneOf) error {
	if expected == nil && oneOf == nil {
		return nil
	}

	if expected == nil && oneOf != nil {
		return errors.New("expected to be nil")
	}

	if expected != nil && oneOf == nil {
		return fmt.Errorf("expected not to be nil")
	}

	if actual, expected := len(oneOf), len(expected); actual != expected {
		return fmt.Errorf("number of elements does not match, got %d while expected %d", actual, expected)
	}

	for index, property := range oneOf {
		if err := property.Compare(expected[index]); err != nil {
			return fmt.Errorf("'oneof' choice mismatch: %w", err)
		}
	}

	return nil
}
