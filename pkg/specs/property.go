package specs

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/jexia/maestro/pkg/specs/labels"
	"github.com/jexia/maestro/pkg/specs/types"
)

// PropertyReference represents a mustach template reference
type PropertyReference struct {
	Resource string    `json:"resource'"`
	Path     string    `json:"path"`
	Property *Property `json:"property"`
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

// Enum represents a enum configuration
type Enum struct {
	Name        string                `json:"name"`
	Values      map[string]*EnumValue `json:"values"`
	Description string                `json:"description"`
}

// EnumValue represents a enum configuration
type EnumValue struct {
	Position    int32  `json:"position"`
	Description string `json:"description"`
}

// ParameterMap is the initial map of parameter names (keys) and their (templated) values (values)
type ParameterMap struct {
	Schema   string               `json:"schema"`
	Options  Options              `json:"options"`
	Header   Header               `json:"header"`
	Property *Property            `json:"property"`
	Stack    map[string]*Property `json:"stack"`
}
