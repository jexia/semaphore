package specs

import (
	"strings"

	"github.com/jexia/maestro/specs/types"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

// IsType checks whether the given value is a type definition
func IsType(value string) bool {
	return strings.HasPrefix(value, types.TypeOpen) && strings.HasSuffix(value, types.TypeClose)
}

// GetTypeContent trims the opening and closing tags from the given type value
func GetTypeContent(value string) string {
	value = strings.Replace(value, types.TypeOpen, "", 1)
	value = strings.Replace(value, types.TypeClose, "", 1)
	value = strings.TrimSpace(value)
	return value
}

// SetType parses the given type and sets the property type
func SetType(property *Property, value cty.Value) {
	if value.Type() != cty.String {
		return
	}

	property.Type = types.Type(GetTypeContent(value.AsString()))
}

// SetDefaultValue sets the given value as default value inside the given property
func SetDefaultValue(property *Property, value cty.Value) {
	switch value.Type() {
	case cty.String:
		property.Default = value.AsString()
		property.Type = types.TypeString
	case cty.Number:
		var def int64
		gocty.FromCtyValue(value, &def)

		property.Default = def
		property.Type = types.TypeInt64
	case cty.Bool:
		var def bool
		gocty.FromCtyValue(value, &def)

		property.Default = def
		property.Type = types.TypeBool
	}
}
