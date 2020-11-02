package specs

import (
	"github.com/jexia/semaphore/pkg/specs/metadata"
	"github.com/jexia/semaphore/pkg/specs/types"
)

// Template contains property schema. This is a union type (Only one field must be set).
type Template struct {
	*metadata.Meta

	// Unique template identifier
	Identifier string

	Reference *PropertyReference `json:"reference,omitempty"` // Reference represents a property reference made inside the given property

	// Only one of the following fields should be set
	Scalar   *Scalar  `json:"scalar,omitempty" yaml:"scalar,omitempty"`
	Enum     *Enum    `json:"enum,omitempty" yaml:"enum,omitempty"`
	Repeated Repeated `json:"repeated,omitempty" yaml:"repeated,omitempty"`
	Message  Message  `json:"message,omitempty" yaml:"message,omitempty"`
	OneOf    OneOf    `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
}

// Type returns the type of the given template.
func (template *Template) Type() types.Type {
	switch {
	case template == nil:
		return types.Unknown
	case template.Message != nil:
		return types.Message
	case template.Repeated != nil:
		return types.Array
	case template.Enum != nil:
		return types.Enum
	case template.Scalar != nil:
		return template.Scalar.Type
	case template.OneOf != nil:
		return types.OneOf
	default:
		return types.Unknown
	}
}

// Clone internal value.
func (template Template) Clone() *Template {
	return template.clone(make(map[string]*Template))
}

func (template *Template) clone(seen map[string]*Template) *Template {
	if template == nil {
		return &Template{}
	}

	var clone = &Template{
		Identifier: template.Identifier,
		Reference:  template.Reference.Clone(),
	}

	if template.Identifier != "" {
		existing, ok := seen[template.Identifier]
		if ok {
			return existing
		}

		seen[template.Identifier] = clone
	}

	if template.Scalar != nil {
		clone.Scalar = template.Scalar.Clone()
	}

	if template.Enum != nil {
		clone.Enum = template.Enum.Clone()
	}

	if template.Message != nil {
		clone.Message = template.Message.clone(seen)
	}

	if template.Repeated != nil {
		clone.Repeated = template.Repeated.clone(seen)
	}

	return clone
}

// Compare given template against the provided one returning the frst mismatch.
func (template Template) Compare(expected *Template) (err error) {
	return template.compare(make(map[string]*Template), expected)
}

func (template *Template) compare(seen map[string]*Template, expected *Template) (err error) {
	if template.Identifier != "" {
		if _, ok := seen[template.Identifier]; ok {
			return nil
		}

		seen[template.Identifier] = template
	}

	switch {
	case template != nil && expected == nil:
		err = errNilTemplate
	case expected.Repeated != nil:
		err = template.Repeated.compare(seen, expected.Repeated)
	case expected.Scalar != nil:
		err = template.Scalar.Compare(expected.Scalar)
	case expected.Message != nil:
		err = template.Message.compare(seen, expected.Message)
	case expected.Enum != nil:
		err = template.Enum.Compare(expected.Enum)
	}

	if err != nil {
		return errTypeMismatch{err}
	}

	return nil
}

// Define ensures that all missing nested template are defined
func (template *Template) Define(expected *Template) {
	template.define(make(map[string]*Template), expected)
}

func (template *Template) define(defined map[string]*Template, expected *Template) {
	if template.Message != nil && expected.Message != nil {
		for key, value := range expected.Message {
			existing, has := template.Message[key]
			if has {
				existing.define(defined, value)
				continue
			}

			template.Message[key] = value.clone(defined)
		}
	}

	// TODO: figure out on how to define repeated
	// this implementation requires that the positions inside the schema and flow
	// are overlapping.

	if template.Message == nil && expected.Message != nil {
		template.Message = expected.Message.clone(defined)
	}

	if template.Enum == nil && expected.Enum != nil {
		template.Enum = expected.Enum.Clone()
	}

	if template.Scalar == nil && expected.Scalar != nil {
		template.Scalar = expected.Scalar.Clone()
	}
}
