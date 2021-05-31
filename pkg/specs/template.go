package specs

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jexia/semaphore/v2/pkg/broker"
	"github.com/jexia/semaphore/v2/pkg/broker/logger"
	"github.com/jexia/semaphore/v2/pkg/specs/metadata"
	"github.com/jexia/semaphore/v2/pkg/specs/types"
	"go.uber.org/zap"
)

var (
	// referencePattern is the matching pattern for references
	referencePattern = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]*:[a-zA-Z0-9\^\&\%\$@_\-\.]*$`)
	// stringPattern is the matching pattern for strings
	stringPattern = regexp.MustCompile(`^\'(.+)\'$`)
)

const (
	// TemplateOpen tag
	TemplateOpen = "{{"
	// TemplateClose tag
	TemplateClose = "}}"

	// PathDelimiter represents the path reference delimiter
	PathDelimiter = "."

	// InputResource key
	InputResource = "input"
	// OutputResource key
	OutputResource = "output"
	// StackResource property
	StackResource = "stack"
	// ParamsResource property
	ParamsResource = "params"
	// RequestResource property
	RequestResource = "request"
	// HeaderResource property
	HeaderResource = "header"
	// ResponseResource property
	ResponseResource = "response"
	// ErrorResource property
	ErrorResource = "error"

	// DefaultInputResource represents the default input property on resource select
	DefaultInputResource = RequestResource
	// DefaultCallResource represents the default call property on resource select
	DefaultCallResource = ResponseResource
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
	OneOf    OneOf    `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
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

	if template.OneOf != nil {
		return types.OneOf
	}

	return types.Unknown
}

// Return the default value for the template (type), assuming the current type has a default.
func (template Template) DefaultValue() interface{} {
	switch {
	case template.Scalar != nil:
		return template.Scalar.Default
	case template.Message != nil:
		return nil
	case template.Repeated != nil:
		return nil
	case template.Enum != nil:
		return nil
	}

	return nil
}

// Clone internal value.
func (template Template) Clone() Template {
	clone := Template{
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

// ShallowClone clones the given template but ignores any nested templates
func (template Template) ShallowClone() Template {
	clone := Template{
		Reference: template.Reference.Clone(),
	}

	if template.Scalar != nil {
		clone.Scalar = template.Scalar.Clone()
	}

	if template.Enum != nil {
		clone.Enum = template.Enum.Clone()
	}

	return clone
}

// Compare given template against the provided one returning the first mismatch.
func (template Template) Compare(expected Template) (err error) {
	switch {
	case expected.Repeated != nil:
		err = template.Repeated.Compare(expected.Repeated)

	case expected.Scalar != nil:
		err = template.Scalar.Compare(expected.Scalar)

	case expected.Message != nil:
		err = template.Message.Compare(expected.Message)

	case expected.Enum != nil:
		err = template.Enum.Compare(expected.Enum)
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

// IsTemplate checks whether the given value is a template
func IsTemplate(value string) bool {
	return strings.HasPrefix(value, TemplateOpen) && strings.HasSuffix(value, TemplateClose)
}

// IsReference returns true if the given value is a reference value
func IsReference(value string) bool {
	return referencePattern.MatchString(value)
}

// IsString returns true if the given value is a string value
func IsString(value string) bool {
	return stringPattern.MatchString(value)
}

// GetTemplateContent trims the opening and closing tags from the given template value
func GetTemplateContent(value string) string {
	value = strings.Replace(value, TemplateOpen, "", 1)
	value = strings.Replace(value, TemplateClose, "", 1)
	value = strings.TrimSpace(value)
	return value
}

// ParseTemplateReference parses the given value as a template reference.
func ParseTemplateReference(value string) (Template, error) {
	if strings.Count(value, "..") > 0 {
		return Template{}, ErrPathNotFound{
			Path: value,
		}
	}

	result := Template{
		Reference: ParsePropertyReference(value),
	}

	return result, nil
}

// ParseTemplateContent parses the given template function.
func ParseTemplateContent(content string) (Template, error) {
	if IsReference(content) {
		return ParseTemplateReference(content)
	}

	if IsString(content) {
		matched := stringPattern.FindStringSubmatch(content)
		result := Template{
			Scalar: &Scalar{
				Type:    types.String,
				Default: matched[1],
			},
		}

		return result, nil
	}

	return Template{}, nil
}

// ParseTemplate parses the given value template, note that path is only used for debugging.
func ParseTemplate(ctx *broker.Context, path string, value string) (Template, error) {
	content := GetTemplateContent(value)
	logger.Debug(ctx, "parsing value template", zap.String("path", path), zap.String("content", content))

	result, err := ParseTemplateContent(content)
	if err != nil {
		return result, err
	}

	logger.Debug(ctx, "parsed template results",
		zap.String("path", path),
		zap.Any("default", result.DefaultValue()),
		zap.String("reference", result.Reference.String()),
	)

	return result, nil
}
