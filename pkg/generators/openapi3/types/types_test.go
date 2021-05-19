package types

import (
	"testing"

	"github.com/jexia/semaphore/v2/pkg/specs/types"
)

func TestOpenTypes(t *testing.T) {
	t.Parallel()

	tests := map[types.Type]Type{
		types.Double:   Integer,
		types.Int64:    Integer,
		types.Uint64:   Integer,
		types.Int32:    Integer,
		types.Uint32:   Integer,
		types.Fixed32:  Integer,
		types.Fixed64:  Integer,
		types.Float:    Integer,
		types.String:   String,
		types.Enum:     String,
		types.Bool:     Boolean,
		types.Bytes:    String,
		types.Sfixed32: Integer,
		types.Sfixed64: Integer,
		types.Sint32:   Integer,
		types.Sint64:   Integer,
	}

	for input, expected := range tests {
		t.Run(string(input), func(t *testing.T) {
			returns := Open(input)
			if returns != expected {
				t.Fatalf("unexpected returned type (%+v), expected (%+v)", returns, expected)
			}
		})
	}
}
