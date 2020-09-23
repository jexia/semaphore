package template

import (
	"regexp"
	"strings"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/broker/trace"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
	"go.uber.org/zap"
)

var (
	// ReferencePattern is the matching pattern for references
	ReferencePattern = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]*:[a-zA-Z0-9\^\&\%\$@_\-\.]*$`)
	// StringPattern is the matching pattern for strings
	StringPattern = regexp.MustCompile(`^\'(.+)\'$`)
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
func ParseReference(path string, name string, value string) (*specs.Property, error) {
	// TODO: check values
	if strings.Count(value, "..") > 0 {
		return nil, trace.New(trace.WithMessage("invalid path, path cannot contain two dots"))
	}

	prop := &specs.Property{
		Name:      name,
		Path:      JoinPath(path, name),
		Reference: ParsePropertyReference(value),
		Raw:       value,
	}

	return prop, nil
}

// ParseContent parses the given template function
func ParseContent(path string, name string, content string) (*specs.Property, error) {
	if ReferencePattern.MatchString(content) {
		return ParseReference(path, name, content)
	}

	if StringPattern.MatchString(content) {
		matched := StringPattern.FindStringSubmatch(content)
		result := &specs.Property{
			Name:  name,
			Path:  path,
			Label: labels.Optional,
			Template: specs.Template{
				Scalar: &specs.Scalar{
					Type:    types.String,
					Default: matched[1],
				},
			},
		}

		return result, nil
	}

	return &specs.Property{
		Name: name,
		Path: path,
		Raw:  content,
	}, nil
}

// Parse parses the given value template and sets the resource and path
func Parse(ctx *broker.Context, path string, name string, value string) (*specs.Property, error) {
	content := GetTemplateContent(value)
	logger.Debug(ctx, "parsing property template", zap.String("path", path), zap.String("template", content))

	result, err := ParseContent(path, name, content)
	if err != nil {
		return nil, err
	}

	logger.Debug(ctx, "template results in property with type",
		zap.String("path", path),
		zap.Any("default", result.Scalar.Default),
		zap.String("reference", result.Reference.String()),
	)

	return result, nil
}
