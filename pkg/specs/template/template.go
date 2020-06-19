package template

import (
	"regexp"
	"strings"

	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/logger"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/sirupsen/logrus"
)

var (
	// ReferencePattern is the matching pattern for references
	ReferencePattern = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]*:[a-zA-Z0-9_\-\.]*$`)
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
	// ResourceParams property
	ResourceParams = "params"
	// ResourceRequest property
	ResourceRequest = "request"
	// ResourceHeader property
	ResourceHeader = "header"
	// ResourceResponse property
	ResourceResponse = "response"

	// DefaultInputProperty represents the default input property on resource select
	DefaultInputProperty = ResourceRequest
	// DefaultCallProperty represents the default call property on resource select
	DefaultCallProperty = ResourceResponse
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

	if prop == ResourceHeader {
		path = strings.ToLower(path)
	}

	reference.Path = path
	return reference
}

// ParseReference parses the given value as a template reference
func ParseReference(path string, name string, value string) *specs.Property {
	prop := &specs.Property{
		Name:      name,
		Path:      JoinPath(path, name),
		Reference: ParsePropertyReference(value),
		Raw:       value,
	}

	return prop
}

// ParseContent parses the given template function
func ParseContent(path string, name string, content string) (*specs.Property, error) {
	if ReferencePattern.MatchString(content) {
		return ParseReference(path, name, content), nil
	}

	return &specs.Property{
		Name: name,
		Path: path,
		Raw:  content,
	}, nil
}

// Parse parses the given value template and sets the resource and path
func Parse(ctx instance.Context, path string, name string, value string) (*specs.Property, error) {
	content := GetTemplateContent(value)
	ctx.Logger(logger.Core).WithField("path", path).WithField("template", content).Debug("Parsing property template")

	result, err := ParseContent(path, name, content)
	if err != nil {
		return nil, err
	}

	ctx.Logger(logger.Core).WithFields(logrus.Fields{
		"path":      path,
		"type":      result.Type,
		"default":   result.Default,
		"reference": result.Reference,
	}).Debug("Template results in property with type")

	return result, nil
}
