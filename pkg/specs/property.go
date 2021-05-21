package specs

import (
	"fmt"

	"github.com/jexia/semaphore/v2/pkg/specs/labels"
	"github.com/jexia/semaphore/v2/pkg/specs/metadata"
	"github.com/jexia/semaphore/v2/pkg/specs/types"
)

// Schemas represents a map string collection of properties
type Schemas map[string]*Property

// Get attempts to return the given key from the objects collection
func (objects Schemas) Get(key string) *Property {
	return objects[key]
}

// Append appends the given objects to the objects collection
func (objects Schemas) Append(arg Schemas) {
	for key, val := range arg {
		objects[key] = val
	}
}

// Expression provides information about expression.
type Expression interface {
	Position() string
}

// Property represents a value property.
type Property struct {
	*metadata.Meta
	Name        string `json:"name,omitempty" yaml:"name,omitempty"`               // Name represents the name of the given property
	Path        string `json:"path,omitempty" yaml:"path,omitempty"`               // Path represents the full path to the given property
	Description string `json:"description,omitempty" yaml:"description,omitempty"` // Description holds the description of the given property used to describe its use

	Position int32 `json:"position,omitempty" yaml:"position,omitempty"` // Position of the given property (in array/object)

	Options Options    `json:"options,omitempty" yaml:"options,omitempty"` // Options holds variable options used inside single modules or components
	Expr    Expression `json:"expression,omitempty"`                       // Expr represents the position on where the given property is defined
	Raw     string     `json:"raw,omitempty"`                              // Raw holds the raw template string used to define the given property

	Label labels.Label `json:"label,omitempty" yaml:"label,omitempty"` // Label label describes the usage of a given property ex: optional

	Template `json:"template" yaml:"template"`
}

// DefaultValue returns the default value for a given property.
func (property *Property) DefaultValue() interface{} {
	return property.Template.DefaultValue()
}

// Empty checks if the property has any defined type
func (property *Property) Empty() bool {
	return property.Type() == types.Unknown
}

// Clone makes a deep clone of the given property
func (property *Property) Clone() *Property {
	if property == nil {
		return &Property{}
	}

	result := property.ShallowClone()
	result.Template = property.Template.Clone()
	return result
}

// ShallowClone clones the given property but ignores the defined template and/or
// nested properties. This method is often used in cases where comparisons between
// the flow and schema are made and any defined properties are seen as defined values.
func (property *Property) ShallowClone() *Property {
	if property == nil {
		return &Property{}
	}

	return &Property{
		Meta:        property.Meta,
		Position:    property.Position,
		Description: property.Description,
		Name:        property.Name,
		Path:        property.Path,

		Expr:     property.Expr,
		Raw:      property.Raw,
		Options:  property.Options,
		Label:    property.Label,
		Template: property.Template.ShallowClone(),
	}
}

// Compare checks the given property against the provided one.
func (property *Property) Compare(expected *Property) error {
	if expected == nil {
		return fmt.Errorf("unable to check types for '%s' no schema given", property.Path)
	}

	if property.Type() != expected.Type() {
		return fmt.Errorf("cannot use type (%s) for '%s', expected (%s)", property.Type(), property.Path, expected.Type())
	}

	if property.Label != expected.Label {
		return fmt.Errorf("cannot use label (%s) for '%s', expected (%s)", property.Label, property.Path, expected.Label)
	}

	if !property.Empty() && expected.Empty() {
		return fmt.Errorf("property '%s' has a nested object but schema does not '%s'", property.Path, expected.Name)
	}

	if !expected.Empty() && property.Empty() {
		return fmt.Errorf("schema '%s' has a nested object but property does not '%s'", expected.Name, property.Path)
	}

	if err := property.Template.Compare(expected.Template); err != nil {
		return fmt.Errorf("nested schema mismatch under property '%s': %w", property.Path, err)
	}

	return nil
}

// Define ensures that all missing nested properties are defined
func (property *Property) Define(expected *Property) {
	property.Position = expected.Position
	property.Template.Define(expected.Template)
}

// ParameterMap is the initial map of parameter names (keys) and their (templated) values (values)
type ParameterMap struct {
	*metadata.Meta
	DependsOn Dependencies         `json:"depends_on,omitempty"`
	Schema    string               `json:"schema,omitempty"`
	Params    map[string]*Property `json:"params,omitempty"`
	Options   Options              `json:"options,omitempty"`
	Header    Header               `json:"header,omitempty"`
	Property  *Property            `json:"property,omitempty"`
	Stack     map[string]*Property `json:"stack,omitempty"`
}

// Clone clones the given parameter map
func (parameters *ParameterMap) Clone() *ParameterMap {
	if parameters == nil {
		return nil
	}

	result := &ParameterMap{
		Meta:     parameters.Meta,
		Schema:   parameters.Schema,
		Params:   make(map[string]*Property, len(parameters.Params)),
		Options:  make(Options, len(parameters.Options)),
		Header:   make(Header, len(parameters.Header)),
		Stack:    make(map[string]*Property, len(parameters.Stack)),
		Property: parameters.Property.Clone(),
	}

	for key, property := range parameters.Params {
		result.Params[key] = property.Clone()
	}

	for key, value := range parameters.Options {
		result.Options[key] = value
	}

	for key, value := range parameters.Header {
		result.Header[key] = value.Clone()
	}

	for key, property := range parameters.Stack {
		result.Stack[key] = property.Clone()
	}

	return result
}
