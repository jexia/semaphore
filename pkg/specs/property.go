package specs

import (
	"encoding/json"

	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/metadata"
	"github.com/jexia/semaphore/pkg/specs/types"
)

// PropertyReference represents a mustach template reference
type PropertyReference struct {
	*metadata.Meta
	Resource string    `json:"resource,omitempty"`
	Path     string    `json:"path,omitempty"`
	Property *Property `json:"-"`
}

// Clone clones the given property reference
func (reference *PropertyReference) Clone() *PropertyReference {
	if reference == nil {
		return nil
	}

	return &PropertyReference{
		Meta:     reference.Meta,
		Resource: reference.Resource,
		Path:     reference.Path,
		Property: nil,
	}
}

func (reference *PropertyReference) String() string {
	if reference == nil {
		return ""
	}

	return reference.Resource + ":" + reference.Path
}

// Schemas represents a map string collection of properties
type Schemas map[string]*Property

// Get attempts to return the given key from the objects collection
func (objects Schemas) Get(key string) *Property {
	if objects == nil {
		return nil
	}

	return objects[key]
}

// Append appends the given objects to the objects collection
func (objects Schemas) Append(arg Schemas) {
	if objects == nil {
		return
	}

	for key, val := range arg {
		objects[key] = val
	}
}

// Expression provides information about expression.
type Expression interface {
	Position() string
}

// Property represents a value property.
// A value property could contain a constant value or a value reference.
type Property struct {
	*metadata.Meta
	Position  int32                `json:"position,omitempty"`
	Comment   string               `json:"comment,omitempty"`
	Name      string               `json:"name,omitempty"`
	Path      string               `json:"path,omitempty"`
	Default   interface{}          `json:"default,omitempty"`
	Type      types.Type           `json:"type,omitempty"`
	Label     labels.Label         `json:"label,omitempty"`
	Reference *PropertyReference   `json:"reference,omitempty"`
	Nested    map[string]*Property `json:"nested,omitempty"`
	Expr      Expression           `json:"-"`
	Raw       string               `json:"raw,omitempty"`
	Options   Options              `json:"options,omitempty"`
	Enum      *Enum                `json:"enum,omitempty"`
}

// UnmarshalJSON corrects the 64bit data types in accordance with golang
func (prop *Property) UnmarshalJSON(data []byte) error {
	type property Property
	p := property{}
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	*prop = Property(p)
	prop.Clean()

	for _, nested := range prop.Nested {
		nested.Clean()
	}

	return nil
}

// Clean fixes the type casting issue of unmarshal
func (prop *Property) Clean() {
	if prop.Default != nil {
		switch prop.Type {
		case types.Int64, types.Sint64, types.Sfixed64:
			_, ok := prop.Default.(int64)
			if !ok {
				prop.Default = int64(prop.Default.(float64))
			}
		case types.Uint64, types.Fixed64:
			_, ok := prop.Default.(uint64)
			if !ok {
				prop.Default = uint64(prop.Default.(float64))
			}
		case types.Int32, types.Sint32, types.Sfixed32:
			_, ok := prop.Default.(int32)
			if !ok {
				prop.Default = int32(prop.Default.(float64))
			}
		case types.Uint32, types.Fixed32:
			_, ok := prop.Default.(uint32)
			if !ok {
				prop.Default = uint32(prop.Default.(float64))
			}
		}
	}
}

// Clone makes a deep clone of the given property
func (prop *Property) Clone() *Property {
	if prop == nil {
		return &Property{}
	}

	result := &Property{
		Meta:      prop.Meta,
		Position:  prop.Position,
		Reference: prop.Reference.Clone(),
		Comment:   prop.Comment,
		Name:      prop.Name,
		Path:      prop.Path,
		Default:   prop.Default,
		Type:      prop.Type,
		Label:     prop.Label,
		Expr:      prop.Expr,
		Raw:       prop.Raw,
		Options:   prop.Options,
		Enum:      prop.Enum,
		Nested:    make(map[string]*Property, len(prop.Nested)),
	}

	for key, nested := range prop.Nested {
		result.Nested[key] = nested.Clone()
	}

	return result
}

// Enum represents a enum configuration
type Enum struct {
	*metadata.Meta
	Name        string                `json:"name,omitempty"`
	Keys        map[string]*EnumValue `json:"keys,omitempty"`
	Positions   map[int32]*EnumValue  `json:"positions,omitempty"`
	Description string                `json:"description,omitempty"`
}

// EnumValue represents a enum configuration
type EnumValue struct {
	*metadata.Meta
	Key         string `json:"key,omitempty"`
	Position    int32  `json:"position,omitempty"`
	Description string `json:"description,omitempty"`
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

	for key, prop := range parameters.Params {
		result.Params[key] = prop.Clone()
	}

	for key, value := range parameters.Options {
		result.Options[key] = value
	}

	for key, value := range parameters.Header {
		result.Header[key] = value.Clone()
	}

	for key, prop := range parameters.Stack {
		result.Stack[key] = prop.Clone()
	}

	return result
}
