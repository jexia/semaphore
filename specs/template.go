package specs

import (
	"strings"

	"github.com/zclconf/go-cty/cty"
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
	// ResourceRequest property
	ResourceRequest = "request"
	// ResourceRequestHeader property
	ResourceRequestHeader = "request.header"
	// ResourceResponse property
	ResourceResponse = "response"
	// ResourceResponseHeader property
	ResourceResponseHeader = "response.header"

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

// ParseReference parses the given value as a template reference
func ParseReference(value string) *PropertyReference {
	rv := strings.Split(value, ReferenceDelimiter)

	reference := PropertyReference{
		Resource: rv[0],
		Label:    LabelOptional,
	}

	if len(rv) == 1 {
		return &reference
	}

	reference.Path = rv[1]
	return &reference
}

// SetTemplate parses the given value template and sets the resource and path
func SetTemplate(property *Property, value cty.Value) {
	if value.Type() != cty.String {
		return
	}

	content := GetTemplateContent(value.AsString())

	// Currently only references could be defined inside templates, possible features could be added in the future
	property.Reference = ParseReference(content)
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
