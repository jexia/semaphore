package types

import (
	"errors"
	"testing"
)

func TestDecodeFromString(t *testing.T) {
	type test struct {
		input    string
		expected interface{}
		error    error
	}

	tests := map[Type]test{
		"unknown": {
			input:    `value`,
			expected: nil,
			error:    ErrUnknownType("unknown"),
		},
		Double: {
			input:    `10`,
			expected: float64(10),
		},
		Int64: {
			input:    `10`,
			expected: int64(10),
		},
		Uint64: {
			input:    `10`,
			expected: uint64(10),
		},
		Fixed64: {
			input:    `10`,
			expected: uint64(10),
		},
		Int32: {
			input:    `10`,
			expected: int32(10),
		},
		Uint32: {
			input:    `10`,
			expected: uint32(10),
		},
		Fixed32: {
			input:    `10`,
			expected: uint32(10),
		},
		Float: {
			input:    `10`,
			expected: float64(10),
		},
		String: {
			input:    `msg`,
			expected: "msg",
		},
		Bool: {
			input:    `true`,
			expected: true,
		},
		Bytes: {
			input:    "aGVsbG83",
			expected: []byte{104, 101, 108, 108, 111, 55},
		},
		Sfixed32: {
			input:    `10`,
			expected: int32(10),
		},
		Sfixed64: {
			input:    `10`,
			expected: int64(10),
		},
		Sint32: {
			input:    `10`,
			expected: int32(10),
		},
		Sint64: {
			input:    `10`,
			expected: int64(10),
		},
	}

	for typed, test := range tests {
		t.Run(string(typed), func(t *testing.T) {
			val, err := DecodeFromString(test.input, typed)

			if !errors.Is(err, test.error) {
				t.Fatal(err)
			}

			switch val.(type) {
			case []byte:
				bb := val.([]byte)
				expected := test.expected.([]byte)

				if len(bb) != len(expected) {
					t.Errorf("unexpected result %+v, expected %+v", bb, expected)
				}
			default:
				if val != test.expected {
					t.Errorf("unexpected result %+v, expected %+v", val, test.expected)
				}
			}
		})
	}
}
