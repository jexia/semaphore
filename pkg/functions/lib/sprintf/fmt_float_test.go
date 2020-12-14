package sprintf

import (
	"errors"
	"testing"

	"github.com/jexia/semaphore/pkg/specs/types"
)

func TestFloatCanFormat(t *testing.T) {
	var tests = map[types.Type]bool{
		types.Bool:     false,
		types.Bytes:    false,
		types.Double:   false,
		types.Enum:     false,
		types.Fixed32:  false,
		types.Fixed64:  false,
		types.Float:    true,
		types.Int32:    false,
		types.Int64:    false,
		types.Message:  false,
		types.Sfixed32: false,
		types.Sfixed64: false,
		types.Sint32:   false,
		types.Sint64:   false,
		types.String:   false,
		types.Uint32:   false,
		types.Uint64:   false,
	}

	var constructor Float

	for dataType, expected := range tests {
		t.Run(string(dataType), func(t *testing.T) {
			if actual := constructor.CanFormat(dataType); actual != expected {
				t.Errorf("unexpected '%v', expected '%v'", actual, expected)
			}
		})
	}
}

func TestFtoa(t *testing.T) {
	type test struct {
		value     interface{}
		precision Precision
		expected  string
		error     error
	}

	var tests = map[string]test{
		"nil value":          {error: errNoValue},
		"not a float type":   {value: int(42), error: errNonFloatType},
		"float32 with scale": {value: float32(3.14159265), precision: Precision{Scale: 2}, expected: "3.14"},
		"float64 with width": {value: float32(3.14159265), precision: Precision{Width: 7}, expected: "3.141593"},
	}

	for title, test := range tests {
		t.Run(title, func(t *testing.T) {
			actual, err := ftoa(test.precision, test.value)

			if !errors.Is(err, test.error) {
				t.Errorf("unexpected error '%s', expected '%s'", err, test.error)
			}

			if actual != test.expected {
				t.Errorf("the output '%s' was expected to be '%s'", actual, test.expected)
			}
		})
	}
}
