package sprintf

import (
	"testing"

	"github.com/jexia/semaphore/pkg/specs/types"
)

func TestJSONCanFormat(t *testing.T) {
	var tests = map[types.Type]bool{
		types.Bool:     true,
		types.Bytes:    true,
		types.Double:   true,
		types.Enum:     true,
		types.Fixed32:  true,
		types.Fixed64:  true,
		types.Float:    true,
		types.Int32:    true,
		types.Int64:    true,
		types.Message:  true,
		types.Sfixed32: true,
		types.Sfixed64: true,
		types.Sint32:   true,
		types.Sint64:   true,
		types.String:   true,
		types.Uint32:   true,
		types.Uint64:   true,
	}

	var constructor JSON

	for dataType, expected := range tests {
		t.Run(string(dataType), func(t *testing.T) {
			if actual := constructor.CanFormat(dataType); actual != expected {
				t.Errorf("unexpected '%v', expected '%v'", actual, expected)
			}
		})
	}
}
