package specs

import (
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
	Default   interface{}        `json:"default,omitempty"`
	Type      types.Type         `json:"type,omitempty"`
	Label     labels.Label       `json:"label,omitempty"`
	Reference *PropertyReference `json:"reference,omitempty"`
}

// Clone scalar value.
func (s Scalar) Clone() *Scalar {
	return &Scalar{
		Default:   s.Default,
		Type:      s.Type,
		Label:     s.Label,
		Reference: s.Reference.Clone(),
	}
}

// Enum represents a enum configuration
type Enum struct {
	*metadata.Meta
	Name        string                `json:"name,omitempty"`
	Keys        map[string]*EnumValue `json:"keys,omitempty"`
	Positions   map[int32]*EnumValue  `json:"positions,omitempty"`
	Description string                `json:"description,omitempty"`
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
	// Template contains the type of repeated values
	Template Template

	// Default contains the static values for certain indexes
	Default map[uint]*Property
}

// Clone repeated.
func (e Repeated) Clone() *Repeated {
	var clone = &Repeated{
		Template: *e.Template.Clone(),
		Default:  make(map[uint]*Property, len(e.Default)),
	}

	for index, prop := range e.Default {
		clone.Default[index] = prop
	}

	return clone
}

// Message represents an object which keeps the original order of keys.
type Message struct {
	Keys       []string
	Properties []*Property
}

// Clone the message.
func (m Message) Clone() *Message {
	var clone = &Message{
		Keys:       make([]string, 0, len(m.Keys)),
		Properties: make([]*Property, 0, len(m.Properties)),
	}

	for index := 0; index < len(m.Keys); index++ {
		clone.Keys[index] = m.Keys[index]
		clone.Properties[index] = m.Properties[index].Clone()
	}

	return clone
}

// Template contains property schema. This is a union type (Only one field must be set).
type Template struct {
	Scalar   *Scalar   `json:"scalar,omitempty"`
	Enum     *Enum     `json:"enum,omitempty"`
	Repeated *Repeated `json:"repeated,omitempty"`
	Message  *Message  `json:"message,omitempty"`
	// TODO: OneOf OneOf
}

// Clone internal value.
func (o Template) Clone() *Template {
	var clone = new(Template)

	switch {
	case o.Scalar != nil:
		clone.Scalar = o.Scalar.Clone()

		break
	case o.Enum != nil:
		clone.Enum = o.Enum.Clone()

		break
	case o.Repeated != nil:
		clone.Repeated = o.Repeated.Clone()

		break
	case o.Message != nil:
		clone.Message = o.Message.Clone()

		break
	}

	return clone
}

// Property represents a value property.
type Property struct {
	*metadata.Meta
	Name        string `json:"name,omitempty"`
	Path        string `json:"path,omitempty"`
	Description string `json:"description,omitempty"`

	Position int32 `json:"position,omitempty"` // what is this?

	Options Options    `json:"options,omitempty"`
	Expr    Expression `json:"-"`

	Raw string `json:"raw,omitempty"`

	Template // contains property schema
}

// PropertyList represents a list of properties
type PropertyList []*Property

// Get attempts to return a property inside the given list with the given name
func (nested PropertyList) Get(key string) *Property {
	for _, item := range nested {
		if item == nil {
			continue
		}

		if item.Name == key {
			return item
		}
	}

	return nil
}

// UnmarshalJSON corrects the 64bit data types in accordance with golang
func (prop *Property) UnmarshalJSON(data []byte) error {
	// type property Property
	// p := property{}
	// err := json.Unmarshal(data, &p)
	// if err != nil {
	// 	return err
	// }
	//
	// *prop = Property(p)
	// prop.Clean()
	//
	// for _, nested := range prop.Nested {
	// 	nested.Clean()
	// }

	return nil
}

// Clean fixes the type casting issue of unmarshal
func (prop *Property) Clean() {
	// if prop.Default != nil {
	// 	switch prop.Type {
	// 	case types.Int64, types.Sint64, types.Sfixed64:
	// 		_, ok := prop.Default.(int64)
	// 		if !ok {
	// 			prop.Default = int64(prop.Default.(float64))
	// 		}
	// 	case types.Uint64, types.Fixed64:
	// 		_, ok := prop.Default.(uint64)
	// 		if !ok {
	// 			prop.Default = uint64(prop.Default.(float64))
	// 		}
	// 	case types.Int32, types.Sint32, types.Sfixed32:
	// 		_, ok := prop.Default.(int32)
	// 		if !ok {
	// 			prop.Default = int32(prop.Default.(float64))
	// 		}
	// 	case types.Uint32, types.Fixed32:
	// 		_, ok := prop.Default.(uint32)
	// 		if !ok {
	// 			prop.Default = uint32(prop.Default.(float64))
	// 		}
	// 	}
	// }
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

		Expr:    prop.Expr,
		Raw:     prop.Raw,
		Options: prop.Options,

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
