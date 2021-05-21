package template

import (
	"regexp"
	"strings"

	"github.com/jexia/semaphore/v2/pkg/broker"
	"github.com/jexia/semaphore/v2/pkg/broker/logger"
	"github.com/jexia/semaphore/v2/pkg/specs"
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

	// ReferenceDelimiter represents the value resource reference delimiter
	ReferenceDelimiter = ":"
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

// Is checks whether the given value is a template
func Is(value string) bool {
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

// ParsePropertyReference parses the given value to a property reference
func ParsePropertyReference(value string) *specs.PropertyReference {
	rv := strings.Split(value, ReferenceDelimiter)

	resource := rv[0]

	resources := SplitPath(resource)
	var prop string

	if len(resources) > 1 {
		prop = JoinPath(resources[1:]...)
	}

	reference := &specs.PropertyReference{
		Resource: resource,
	}

	if len(rv) == 1 {
		return reference
	}

	path := rv[1]

	if prop == HeaderResource {
		path = strings.ToLower(path)
	}

	reference.Path = path
	return reference
}

// ParseReference parses the given value as a template reference
func ParseReference(value string) (specs.Template, error) {
	if strings.Count(value, "..") > 0 {
		return specs.Template{}, ErrPathNotFound{
			Path: value,
		}
	}

	result := specs.Template{
		Reference: ParsePropertyReference(value),
	}

	return result, nil
}

// ParseContent parses the given template function
func ParseContent(content string) (specs.Template, error) {
	if IsReference(content) {
		return ParseReference(content)
	}

	if IsString(content) {
		matched := stringPattern.FindStringSubmatch(content)
		result := specs.Template{
			Scalar: &specs.Scalar{
				Type:    types.String,
				Default: matched[1],
			},
		}

		return result, nil
	}

	return specs.Template{}, nil
}

// Parse parses the given value template, note that path is only used for debugging
func Parse(ctx *broker.Context, path string, value string) (specs.Template, error) {
	content := GetTemplateContent(value)
	logger.Debug(ctx, "parsing value template", zap.String("path", path), zap.String("content", content))

	result, err := ParseContent(content)
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
