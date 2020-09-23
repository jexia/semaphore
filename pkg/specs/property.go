package specs

import (
	"encoding/json"
	"sort"

	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/metadata"
	"github.com/jexia/semaphore/pkg/specs/types"
)

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

// Scalar value.
type Scalar struct {
	Default interface{} `json:"default,omitempty"`
	Type    types.Type  `json:"type,omitempty"`
}

// UnmarshalJSON corrects the 64bit data types in accordance with golang
func (scalar *Scalar) UnmarshalJSON(data []byte) error {
	if scalar == nil {
		return nil
	}

	type sc Scalar
	t := sc{}

	err := json.Unmarshal(data, &t)
	if err != nil {
		return err
	}

	*scalar = Scalar(t)
	scalar.Clean()

	return nil
}

// Clean fixes the type casting issue of unmarshal
func (scalar *Scalar) Clean() {
	if scalar.Default == nil {
		return
	}

	switch scalar.Type {
	case types.Int64, types.Sint64, types.Sfixed64:
		_, ok := scalar.Default.(int64)
		if !ok {
			scalar.Default = int64(scalar.Default.(float64))
		}
	case types.Uint64, types.Fixed64:
		_, ok := scalar.Default.(uint64)
		if !ok {
			scalar.Default = uint64(scalar.Default.(float64))
		}
	case types.Int32, types.Sint32, types.Sfixed32:
		_, ok := scalar.Default.(int32)
		if !ok {
			scalar.Default = int32(scalar.Default.(float64))
		}
	case types.Uint32, types.Fixed32:
		_, ok := scalar.Default.(uint32)
		if !ok {
			scalar.Default = uint32(scalar.Default.(float64))
		}
	}
}

// Clone scalar value.
func (scalar Scalar) Clone() *Scalar {
	return &Scalar{
		Default: scalar.Default,
		Type:    scalar.Type,
	}
}

// Enum represents a enum configuration
type Enum struct {
	*metadata.Meta
	Name        string                `json:"name,omitempty"`
	Description string                `json:"description,omitempty"`
	Keys        map[string]*EnumValue `json:"keys,omitempty"`
	Positions   map[int32]*EnumValue  `json:"positions,omitempty"`
}

// Clone enum schema.
func (e Enum) Clone() *Enum { return &e }

// EnumValue represents a enum configuration
type EnumValue struct {
	*metadata.Meta
	Key         string `json:"key,omitempty"`
	Position    int32  `json:"position,omitempty"`
	Description string `json:"description,omitempty"`
}

// Repeated represents an array type.
type Repeated struct {
	*Property

	// Default contains the static values for certain indexes
	Default   map[uint]*Property `json:"default,omitempty"`
	Reference *PropertyReference `json:"reference,omitempty"`
}

// Clone repeated.
func (repeated Repeated) Clone() *Repeated {
	var clone = &Repeated{
		Property: repeated.Property.Clone(),
		Default:  make(map[uint]*Property, len(repeated.Default)),
	}

	for index, prop := range repeated.Default {
		clone.Default[index] = prop
	}

	return clone
}

// Message represents an object which keeps the original order of keys.
type Message map[string]*Property

// SortedProperties returns the available properties as a properties list
// ordered base on the properties position.
func (message Message) SortedProperties() PropertyList {
	result := make(PropertyList, 0, len(message))

	for _, property := range message {
		result = append(result, property)
	}

	sort.Sort(result)
	return result
}

// Clone the message.
func (message Message) Clone() Message {
	var clone = make(map[string]*Property, len(message))

	for key := range message {
		clone[key] = message[key].Clone()
	}

	return clone
}

// Template contains property schema. This is a union type (Only one field must be set).
type Template struct {
	// Only one of the following fields should be set
	Scalar   *Scalar   `json:"scalar,omitempty"`
	Enum     *Enum     `json:"enum,omitempty"`
	Repeated *Repeated `json:"repeated,omitempty"`
	Message  Message   `json:"message,omitempty"`
}

// Type returns the type of the given template.
func (t Template) Type() types.Type {
	if t.Message != nil {
		return types.Message
	}

	if t.Repeated != nil {
		return types.Array
	}

	if t.Enum != nil {
		return types.Enum
	}

	if t.Scalar != nil {
		return t.Scalar.Type
	}

	return types.Unknown
}

// Clone internal value.
func (t Template) Clone() *Template {
	var clone = new(Template)

	if t.Scalar != nil {
		clone.Scalar = t.Scalar.Clone()
	}

	if t.Enum != nil {
		clone.Enum = t.Enum.Clone()
	}

	if t.Repeated != nil {
		clone.Repeated = t.Repeated.Clone()
	}

	if t.Message != nil {
		clone.Message = t.Message.Clone()
	}

	return clone
}

// PropertyList represents a list of properties
type PropertyList []*Property

func (list PropertyList) Len() int           { return len(list) }
func (list PropertyList) Swap(i, j int)      { list[i], list[j] = list[j], list[i] }
func (list PropertyList) Less(i, j int) bool { return list[i].Position < list[j].Position }

// Get attempts to return a property inside the given list with the given name
func (list PropertyList) Get(key string) *Property {
	for _, item := range list {
		if item == nil {
			continue
		}

		if item.Name == key {
			return item
		}
	}

	return nil
}

// Property represents a value property.
type Property struct {
	*metadata.Meta
	Name        string `json:"name,omitempty"`        // Name represents the name of the given property
	Path        string `json:"path,omitempty"`        // Path represents the full path to the given property
	Description string `json:"description,omitempty"` // Description holds the description of the given property used to describe its use

	Position int32 `json:"position,omitempty"` // Position of the given property

	Options Options    `json:"options,omitempty"`    // Options holds variable options used inside single modules or components
	Expr    Expression `json:"expression,omitempty"` // Expr represents the position on where the given property is defined

	Reference *PropertyReference `json:"reference,omitempty"` // Reference represents a property reference made inside the given property
	Raw       string             `json:"raw,omitempty"`       // Raw holds the raw template string used to define the given property

	// Label is the set of field attributes/properties. E.g. label describes that the property is optional.
	Label labels.Label `json:"label,omitempty"`

	Template
}

// Empty checks if the property has any defined type
func (prop *Property) Empty() bool {
	return prop.Template.Type() == types.Unknown
}

// Clone makes a deep clone of the given property
func (prop *Property) Clone() *Property {
	if prop == nil {
		return &Property{}
	}

	return &Property{
		Meta:        prop.Meta,
		Position:    prop.Position,
		Description: prop.Description,
		Name:        prop.Name,
		Path:        prop.Path,

		Expr:      prop.Expr,
		Raw:       prop.Raw,
		Options:   prop.Options,
		Reference: prop.Reference.Clone(),
		Label:     prop.Label,

		Template: *prop.Template.Clone(),
	}
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
