package specs

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/specs/metadata"
	"github.com/jexia/semaphore/pkg/specs/types"
)

// Template contains property schema. This is a union type (Only one field must be set).
type Template struct {
	*metadata.Meta
	Reference *PropertyReference `json:"reference,omitempty"` // Reference represents a property reference made inside the given property

	// Only one of the following fields should be set
	Scalar   *Scalar  `json:"scalar,omitempty" yaml:"scalar,omitempty"`
	Enum     *Enum    `json:"enum,omitempty" yaml:"enum,omitempty"`
	Repeated Repeated `json:"repeated,omitempty" yaml:"repeated,omitempty"`
	Message  Message  `json:"message,omitempty" yaml:"message,omitempty"`
}

// Type returns the type of the given template.
func (template Template) Type() types.Type {
	if template.Message != nil {
		return types.Message
	}

	if template.Repeated != nil {
		return types.Array
	}

	if template.Enum != nil {
		return types.Enum
	}

	if template.Scalar != nil {
		return template.Scalar.Type
	}

	return types.Unknown
}

// Clone internal value.
func (template Template) Clone() Template {
	var clone = Template{
		Reference: template.Reference.Clone(),
	}

	if template.Scalar != nil {
		clone.Scalar = template.Scalar.Clone()
	}

	if template.Enum != nil {
		clone.Enum = template.Enum.Clone()
	}

	if template.Repeated != nil {
		clone.Repeated = template.Repeated.Clone()
	}

	if template.Message != nil {
		clone.Message = template.Message.Clone()
	}

	return clone
}

// Compare given template against the provided one returning the frst mismatch.
func (template Template) Compare(expected Template) (err error) {
	switch {
	case expected.Repeated != nil:
		err = template.Repeated.Compare(expected.Repeated)
		break

	case expected.Scalar != nil:
		err = template.Scalar.Compare(expected.Scalar)
		break

	case expected.Message != nil:
		err = template.Message.Compare(expected.Message)
		break

	case expected.Enum != nil:
		err = template.Enum.Compare(expected.Enum)
		break
	}

	if err != nil {
		return fmt.Errorf("type mismatch: %w", err)
	}

	return nil
}

// Define ensures that all missing nested template are defined
func (template *Template) Define(expected Template) {
	if template.Message != nil && expected.Message != nil {
		for key, value := range expected.Message {
			existing, has := template.Message[key]
			if has {
				existing.Define(value)
				continue
			}

			template.Message[key] = value.Clone()
		}
	}

	// TODO: figure out on how to define repeated
	// this implementation requires that the positions inside the schema and flow
	// are overlapping.

	if template.Message == nil && expected.Message != nil {
		template.Message = expected.Message.Clone()
	}

	if template.Enum == nil && expected.Enum != nil {
		template.Enum = expected.Enum.Clone()
	}

	if template.Scalar == nil && expected.Scalar != nil {
		template.Scalar = expected.Scalar.Clone()
	}
}
