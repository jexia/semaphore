package formencoded

import (
	"testing"

	"github.com/jexia/semaphore/pkg/specs/types"
)

func TestCastType(t *testing.T) {
	type test struct {
		input interface{}
		want  string
	}

	tests := map[types.Type]test{
		types.Double: {
			input: float64(10),
			want:  "1E+01",
		},
		types.Int64: {
			input: int64(10),
			want:  "10",
		},
		types.Uint64: {
			input: uint64(10),
			want:  "10",
		},
		types.Fixed64: {
			input: uint64(10),
			want:  "10",
		},
		types.Int32: {
			input: int32(10),
			want:  "10",
		},
		types.Uint32: {
			input: uint32(10),
			want:  "10",
		},
		types.Fixed32: {
			input: uint32(10),
			want:  "10",
		},
		types.Float: {
			input: float32(10),
			want:  "1E+01",
		},
		types.String: {
			input: "msg",
			want:  "msg",
		},
		types.Bool: {
			input: true,
			want:  "true",
		},
		types.Bytes: {
			input: []byte{10, 10},
			want:  "Cgo=",
		},
		types.Sfixed32: {
			input: int32(10),
			want:  "10",
		},
		types.Sfixed64: {
			input: int64(10),
			want:  "10",
		},
		types.Sint32: {
			input: int32(10),
			want:  "10",
		},
		types.Sint64: {
			input: int64(10),
			want:  "10",
		},
	}

	for typed, test := range tests {
		t.Run(string(typed), func(t *testing.T) {
			got := castType(typed, test.input)

			if got != test.want {
				t.Errorf("unexpected result %s, want %s", got, test.want)
			}
		})
	}
}

func TestAddEmptyType(t *testing.T) {
	tests := map[types.Type]interface{}{
		types.Double:   nil,
		types.Int64:    nil,
		types.Uint64:   nil,
		types.Fixed64:  nil,
		types.Int32:    nil,
		types.Uint32:   nil,
		types.Fixed32:  nil,
		types.Float:    nil,
		types.String:   nil,
		types.Bool:     nil,
		types.Bytes:    nil,
		types.Sfixed32: nil,
		types.Sfixed64: nil,
		types.Sint32:   nil,
		types.Sint64:   nil,
	}

	for typed, value := range tests {
		t.Run(string(typed), func(t *testing.T) {
			got := castType(typed, value)
			want := ""
			if got != want {
				t.Errorf("unexpected result %s", got)
			}
		})
	}
}
