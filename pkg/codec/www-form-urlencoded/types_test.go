package formencoded

import (
	"net/url"
	"testing"

	"github.com/jexia/semaphore/pkg/specs/types"
)

func TestAddType(t *testing.T) {
	type test struct {
		input    interface{}
		expected string
	}

	tests := map[types.Type]test{
		types.Double: {
			input:    float64(10),
			expected: "mock=1E%2B01",
		},
		types.Int64: {
			input:    int64(10),
			expected: "mock=10",
		},
		types.Uint64: {
			input:    uint64(10),
			expected: "mock=10",
		},
		types.Fixed64: {
			input:    uint64(10),
			expected: "mock=10",
		},
		types.Int32: {
			input:    int32(10),
			expected: "mock=10",
		},
		types.Uint32: {
			input:    uint32(10),
			expected: "mock=10",
		},
		types.Fixed32: {
			input:    uint32(10),
			expected: "mock=10",
		},
		types.Float: {
			input:    float32(10),
			expected: "mock=1E%2B01",
		},
		types.String: {
			input:    string("msg"),
			expected: "mock=msg",
		},
		types.Bool: {
			input:    true,
			expected: "mock=true",
		},
		types.Bytes: {
			input:    []byte{10, 10},
			expected: "mock=Cgo%3D",
		},
		types.Sfixed32: {
			input:    int32(10),
			expected: "mock=10",
		},
		types.Sfixed64: {
			input:    int64(10),
			expected: "mock=10",
		},
		types.Sint32: {
			input:    int32(10),
			expected: "mock=10",
		},
		types.Sint64: {
			input:    int64(10),
			expected: "mock=10",
		},
	}

	for typed, test := range tests {
		t.Run(string(typed), func(t *testing.T) {
			encoder := url.Values{}
			AddTypeKey(encoder, "mock", typed, test.input)

			result := encoder.Encode()
			if result == "" {
				t.Errorf("unexpected empty result")
			}

			if result != test.expected {
				t.Errorf("unexpected result %s, expected %s", result, test.expected)
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
			encoder := url.Values{}
			key := "key"
			expected := key + "="

			AddTypeKey(encoder, key, typed, value)

			result := encoder.Encode()
			if result != expected {
				t.Errorf("unexpected result %s", result)
			}
		})
	}
}

func TestDecodeType(t *testing.T) {
	type test struct {
		input    string
		expected interface{}
	}

	tests := map[types.Type]test{
		types.Double: {
			input:    `10`,
			expected: float64(10),
		},
		types.Int64: {
			input:    `10`,
			expected: int64(10),
		},
		types.Uint64: {
			input:    `10`,
			expected: uint64(10),
		},
		types.Fixed64: {
			input:    `10`,
			expected: uint64(10),
		},
		types.Int32: {
			input:    `10`,
			expected: int32(10),
		},
		types.Uint32: {
			input:    `10`,
			expected: uint32(10),
		},
		types.Fixed32: {
			input:    `10`,
			expected: uint32(10),
		},
		types.Float: {
			input:    `10`,
			expected: float64(10),
		},
		types.String: {
			input:    `msg`,
			expected: "msg",
		},
		types.Bool: {
			input:    `true`,
			expected: true,
		},
		types.Bytes: {
			input:    "aGVsbG83",
			expected: []byte{104, 101, 108, 108, 111, 55},
		},
		types.Sfixed32: {
			input:    `10`,
			expected: int32(10),
		},
		types.Sfixed64: {
			input:    `10`,
			expected: int64(10),
		},
		types.Sint32: {
			input:    `10`,
			expected: int32(10),
		},
		types.Sint64: {
			input:    `10`,
			expected: int64(10),
		},
	}

	for typed, test := range tests {
		t.Run(string(typed), func(t *testing.T) {
			val, err := DecodeType(test.input, typed)
			if err != nil {
				t.Fatal(err)
			}

			if val == nil {
				t.Fatal("unexpected nil value")
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
