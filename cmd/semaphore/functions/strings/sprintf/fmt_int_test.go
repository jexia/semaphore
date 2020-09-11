package sprintf

import (
	"errors"
	"testing"

	"github.com/jexia/semaphore/pkg/specs/types"
)

func TestIntCanFormat(t *testing.T) {
	var tests = map[types.Type]bool{
		types.Bool:     false,
		types.Bytes:    false,
		types.Double:   false,
		types.Enum:     false,
		types.Fixed32:  false,
		types.Fixed64:  false,
		types.Float:    false,
		types.Int32:    true,
		types.Int64:    true,
		types.Message:  false,
		types.Sfixed32: false,
		types.Sfixed64: false,
		types.Sint32:   false,
		types.Sint64:   false,
		types.String:   false,
		types.Uint32:   true,
		types.Uint64:   true,
	}

	var constructor Int

	for dataType, expected := range tests {
		t.Run(string(dataType), func(t *testing.T) {
			if actual := constructor.CanFormat(dataType); actual != expected {
				t.Errorf("unexpected '%v', expected '%v'", actual, expected)
			}
		})
	}
}

func TestItoa(t *testing.T) {
	type test struct {
		value     interface{}
		precision Precision
		expected  string
		error     error
	}

	var tests = map[string]test{
		"nil value":           {error: errNoValue},
		"not an integer type": {value: true, error: errNonIntegerType},
		"int":                 {value: int(42), expected: "42"},
		"int32":               {value: int32(42), expected: "42"},
		"int64":               {value: int64(42), expected: "42"},
		"uint":                {value: uint(42), expected: "42"},
		"uint32":              {value: uint32(42), expected: "42"},
		"uint64":              {value: uint64(42), expected: "42"},
	}

	for title, test := range tests {
		t.Run(title, func(t *testing.T) {
			actual, err := itoa(test.precision, test.value)

			if !errors.Is(err, test.error) {
				t.Errorf("unexpected error '%s', expected '%s'", err, test.error)
			}

			if actual != test.expected {
				t.Errorf("the output '%s' was expected to be '%s'", actual, test.expected)
			}
		})
	}
}
