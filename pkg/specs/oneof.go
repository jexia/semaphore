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
type OneOf map[string]*Property

func (oneOf OneOf) String() string { return dump(oneOf) }

// Clone OneOf.
func (oneOf OneOf) Clone() OneOf {
	return OneOf(Message(oneOf).Clone())
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

	for key, property := range oneOf {
		nested, ok := expected[key]
		if !ok {
			return fmt.Errorf("unknown choice '%s'", key)
		}

		if err := property.Compare(nested); err != nil {
			return fmt.Errorf("'oneof' choice mismatch: %w", err)
		}
	}

	return nil
}
