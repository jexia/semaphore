package specs

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/jexia/maestro/pkg/specs/labels"
	"github.com/jexia/maestro/pkg/specs/types"
)

// PropertyReference represents a mustach template reference
type PropertyReference struct {
	Resource string `json:"resource'"`
	Path     string `json:"path"`
	Property *Property
}

// Clone clones the given property reference
func (reference *PropertyReference) Clone() *PropertyReference {
	return &PropertyReference{
		Resource: reference.Resource,
		Path:     reference.Path,
		Property: nil,
	}
}

func (reference *PropertyReference) String() string {
	return reference.Resource + ":" + reference.Path
}

// Property represents a value property.
// A value property could contain a constant value or a value reference.
type Property struct {
	Position  int32                `json:"position"`
	Comment   string               `json:"comment"`
	Name      string               `json:"name"`
	Path      string               `json:"path"`
	Default   interface{}          `json:"default"`
	Type      types.Type           `json:"type"`
	Label     labels.Label         `json:"label"`
	Reference *PropertyReference   `json:"reference"`
	Nested    map[string]*Property `json:"nested"`
	Expr      hcl.Expression       `json:"expr"` // TODO: replace this with a custom solution
	Raw       string               `json:"raw"`
	Options   Options              `json:"options"`
	Enum      *Enum                `json:"enum"`
}

// Clone makes a deep clone of the given property
func (prop *Property) Clone() *Property {
	if prop == nil {
		return &Property{}
	}

	result := &Property{
		Position: prop.Position,
		Comment:  prop.Comment,
		Name:     prop.Name,
		Path:     prop.Path,
		Default:  prop.Default,
		Type:     prop.Type,
		Label:    prop.Label,
		Expr:     prop.Expr,
		Raw:      prop.Raw,
		Options:  prop.Options,
		Enum:     prop.Enum,
		Nested:   make(map[string]*Property, len(prop.Nested)),
	}

	if prop.Reference != nil {
		result.Reference = prop.Reference.Clone()
	}

	for key, nested := range prop.Nested {
		result.Nested[key] = nested.Clone()
	}

	return result
}

// Enum represents a enum configuration
type Enum struct {
	Name        string                `json:"name"`
	Keys        map[string]*EnumValue `json:"keys"`
	Positions   map[int32]*EnumValue  `json:"positions"`
	Description string                `json:"description"`
}

// EnumValue represents a enum configuration
type EnumValue struct {
	Key         string `json:"key"`
	Position    int32  `json:"position"`
	Description string `json:"description"`
}

// ParameterMap is the initial map of parameter names (keys) and their (templated) values (values)
type ParameterMap struct {
	Schema   string               `json:"schema"`
	Params   map[string]*Property `json:"params"`
	Options  Options              `json:"options"`
	Header   Header               `json:"header"`
	Property *Property            `json:"property"`
	Stack    map[string]*Property `json:"stack"`
}

// Clone clones the given parameter map
func (parameters *ParameterMap) Clone() *ParameterMap {
	result := &ParameterMap{
		Schema:   parameters.Schema,
		Params:   make(map[string]*Property, len(parameters.Params)),
		Options:  make(Options, len(parameters.Options)),
		Header:   make(Header, len(parameters.Header)),
		Stack:    make(map[string]*Property, len(parameters.Stack)),
		Property: parameters.Property.Clone(),
	}

	for key, prop := range parameters.Params {
		result.Params[key] = prop.Clone()
	}

	for key, value := range parameters.Options {
		result.Options[key] = value
	}

	for key, value := range parameters.Header {
		result.Header[key] = value
	}

	for key, prop := range parameters.Stack {
		result.Header[key] = prop.Clone()
	}

	return result
}
