package specs

import (
	"regexp"
	"strings"

	"github.com/jexia/maestro/internal/instance"
	"github.com/jexia/maestro/internal/logger"
	"github.com/sirupsen/logrus"
)

var (
	// ReferencePattern is the matching pattern for references
	ReferencePattern = regexp.MustCompile(`^[a-zA-Z0-9_\.]*:[a-zA-Z0-9_\.]*$`)
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

// IsTemplate checks whether the given value is a template
func IsTemplate(value string) bool {
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
func ParsePropertyReference(value string) *PropertyReference {
	rv := strings.Split(value, ReferenceDelimiter)
	reference := &PropertyReference{
		Resource: rv[0],
	}

	if len(rv) == 1 {
		return reference
	}

	reference.Path = rv[1]
	return reference
}

// ParseReference parses the given value as a template reference
func ParseReference(path string, name string, template string) *Property {
	prop := &Property{
		Name:      name,
		Path:      JoinPath(path, name),
		Reference: ParsePropertyReference(template),
		Raw:       template,
	}

	return prop
}

// ParseTemplateContent parses the given template function
func ParseTemplateContent(path string, name string, content string) (*Property, error) {
	if ReferencePattern.MatchString(content) {
		return ParseReference(path, name, content), nil
	}

	return &Property{
		Name: name,
		Path: path,
		Raw:  content,
	}, nil
}

// ParseTemplate parses the given value template and sets the resource and path
func ParseTemplate(ctx instance.Context, path string, name string, value string) (*Property, error) {
	content := GetTemplateContent(value)
	ctx.Logger(logger.Core).WithField("path", path).WithField("template", content).Debug("Parsing property template")

	result, err := ParseTemplateContent(path, name, content)
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

// JoinPath joins the given flow paths
func JoinPath(values ...string) (result string) {
	for _, value := range values {
		if value == "" {
			continue
		}

		if len(result) > 0 && string(result[len(result)-1]) != "." {
			result += "."
		}

		result += value
	}

	if result == "" || result == "." {
		return result
	}

	if string(result[len(result)-1]) == "." {
		result = result[:len(result)-1]
	}

	if string(result[0]) == "." {
		result = result[1:]
	}

	return result
}

// SplitPath splits the given path into parts
func SplitPath(path string) []string {
	return strings.Split(path, PathDelimiter)
}
