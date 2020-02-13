package specs

import (
	"math/big"
	"testing"

	"github.com/jexia/maestro/specs/types"
	"github.com/zclconf/go-cty/cty"
)

func TestIsType(t *testing.T) {
	tests := map[string]bool{
		"<string":  false,
		"string>":  false,
		"string":   false,
		"<string>": true,
		"<int32>":  true,
		"<int64>":  true,
	}

	for input, expected := range tests {
		result := IsType(input)
		if result != expected {
			t.Errorf("unexpected result %t, expected %t %s", result, expected, input)
		}
	}
}

func TestGetTypeContent(t *testing.T) {
	tests := map[string]string{
		"<string>": "string",
		"<int32>":  "int32",
		"int32>":   "int32",
		"<string":  "string",
	}

	for input, expected := range tests {
		result := GetTypeContent(input)
		if expected != result {
			t.Errorf("unexpected result %s, expected %s", result, expected)
		}
	}
}

func TestSetType(t *testing.T) {
	tests := map[string]types.Type{
		"<string>": types.TypeString,
		"<bool>":   types.TypeBool,
		"<int32>":  types.TypeInt32,
		"<int64>":  types.TypeInt64,
	}

	for input, expected := range tests {
		property := Property{}
		SetType(&property, cty.StringVal(input))

		if property.Type != expected {
			t.Errorf("unexpected result %s, expected %s", property.Type, expected)
		}
	}
}

func TestSetDefaultValue(t *testing.T) {
	type expected struct {
		Default interface{}
		Type    types.Type
	}

	tests := map[cty.Value]expected{
		cty.StringVal("default"): {
			Default: "default",
			Type:    types.TypeString,
		},
		cty.NumberVal(big.NewFloat(10)): {
			Default: int64(10),
			Type:    types.TypeInt64,
		},
		cty.BoolVal(true): {
			Default: true,
			Type:    types.TypeBool,
		},
	}

	for input, expected := range tests {
		property := Property{}
		SetDefaultValue(&property, input)

		if expected.Default != property.Default {
			t.Errorf("unexpected result %+v, expected %+v", property.Default, expected.Default)
		}

		if expected.Type != property.Type {
			t.Errorf("unexpected type %s, expected %s", property.Type, expected.Type)
		}
	}
}
