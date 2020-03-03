package specs

import (
	"regexp"
	"strings"

	"github.com/jexia/maestro/specs/trace"
	"github.com/jexia/maestro/specs/types"
	log "github.com/sirupsen/logrus"
)

var (
	// FunctionPattern is the matching pattern for custom defined functions
	FunctionPattern = regexp.MustCompile(`(\w+)\((.*)\)$`)
)

const (
	// TemplateOpen tag
	TemplateOpen = "{{"
	// TemplateClose tag
	TemplateClose = "}}"

	// FunctionArgumentDelimiter represents the character delimiting function arguments
	FunctionArgumentDelimiter = ","
	// ReferenceDelimiter represents the value resource reference delimiter
	ReferenceDelimiter = ":"
	// PathDelimiter represents the path reference delimiter
	PathDelimiter = "."

	// InputResource key
	InputResource = "input"
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
func ParseReference(path string, value string) *Property {
	prop := &Property{
		Path:      path,
		Reference: ParsePropertyReference(value),
	}

	prop.Reference.Label = types.LabelOptional
	return prop
}

// ParseFunction attempts to parses the given function
func ParseFunction(path string, functions CustomDefinedFunctions, content string) (*Property, error) {
	pattern := FunctionPattern.FindStringSubmatch(content)
	name := pattern[1]
	args := strings.Split(pattern[2], FunctionArgumentDelimiter)

	if functions[name] == nil {
		return nil, trace.New(trace.WithMessage("undefined custom function '%s' in '%s'", name, content))
	}

	arguments := make([]*Property, len(args))

	for index, arg := range args {
		result, err := ParseTemplateContent(path, functions, strings.TrimSpace(arg))
		if err != nil {
			return nil, err
		}

		arguments[index] = result
	}

	property, err := functions[name](path, arguments...)
	if err != nil {
		return nil, err
	}

	return property, nil
}

// ParseTemplateContent parses the given template function
func ParseTemplateContent(path string, functions CustomDefinedFunctions, content string) (*Property, error) {
	if FunctionPattern.MatchString(content) {
		return ParseFunction(path, functions, content)
	}

	// TODO: handle constant
	return ParseReference(path, content), nil
}

// ParseTemplate parses the given value template and sets the resource and path
func ParseTemplate(path string, functions CustomDefinedFunctions, value string) (*Property, error) {
	content := GetTemplateContent(value)
	log.WithField("path", path).WithField("template", content).Debug("Parsing property template")

	result, err := ParseTemplateContent(path, functions, content)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
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

	if result == "" {
		return result
	}

	if string(result[len(result)-1]) == "." {
		result = result[:len(result)-1]
	}

	return result
}

// SplitPath splits the given path into parts
func SplitPath(path string) []string {
	return strings.Split(path, PathDelimiter)
}
