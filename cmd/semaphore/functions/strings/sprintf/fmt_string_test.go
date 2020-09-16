package sprintf

import (
	"errors"
	"testing"

	"github.com/jexia/semaphore/pkg/specs/types"
)

func TestStringCanFormat(t *testing.T) {
	var tests = map[types.Type]bool{
		types.Bool:     false,
		types.Bytes:    false,
		types.Double:   false,
		types.Enum:     false,
		types.Fixed32:  false,
		types.Fixed64:  false,
		types.Float:    false,
		types.Int32:    false,
		types.Int64:    false,
		types.Message:  false,
		types.Sfixed32: false,
		types.Sfixed64: false,
		types.Sint32:   false,
		types.Sint64:   false,
		types.String:   true,
		types.Uint32:   false,
		types.Uint64:   false,
	}

	var constructor String

	for dataType, expected := range tests {
		t.Run(string(dataType), func(t *testing.T) {
			if actual := constructor.CanFormat(dataType); actual != expected {
				t.Errorf("unexpected '%v', expected '%v'", actual, expected)
			}
		})
	}
}

func TestStrtoa(t *testing.T) {
	type test struct {
		value     interface{}
		precision Precision
		expected  string
		error     error
	}

	var tests = map[string]test{
		"nil value":        {error: errNoValue},
		"not a float type": {value: true, error: errNonStringType},
		"unlimited string": {value: "this is a string with unlimited length", expected: "this is a string with unlimited length"},
		"limited by width": {value: "this is a string limited by width", precision: Precision{Width: 16}, expected: "this is a string"},
	}

	for title, test := range tests {
		t.Run(title, func(t *testing.T) {
			actual, err := strtoa(test.precision, test.value)

			if !errors.Is(err, test.error) {
				t.Errorf("unexpected error '%s', expected '%s'", err, test.error)
			}

			if actual != test.expected {
				t.Errorf("the output '%s' was expected to be '%s'", actual, test.expected)
			}
		})
	}
}
