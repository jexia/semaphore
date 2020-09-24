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
			input:    "msg",
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
