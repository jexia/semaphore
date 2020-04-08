package json

import (
	"io/ioutil"
	"testing"

	"github.com/francoispqt/gojay"
	"github.com/jexia/maestro/specs/types"
)

type mockObject struct {
	fn func(*gojay.Encoder)
}

func (m *mockObject) MarshalJSONObject(encoder *gojay.Encoder) {
	m.fn(encoder)
}

func (m *mockObject) IsNil() bool {
	return m == nil
}

func TestAddTypeKey(t *testing.T) {
	tests := map[types.Type]interface{}{
		types.Double:   float64(10),
		types.Int64:    int64(10),
		types.Uint64:   uint64(10),
		types.Fixed64:  uint64(10),
		types.Int32:    int32(10),
		types.Uint32:   uint32(10),
		types.Fixed32:  uint64(10),
		types.Float:    float32(10),
		types.String:   string("msg"),
		types.Bool:     true,
		types.Bytes:    []byte{10, 10},
		types.Sfixed32: int32(10),
		types.Sfixed64: int64(10),
		types.Sint32:   int32(10),
		types.Sint64:   int64(10),
	}

	for typed, value := range tests {
		t.Run(string(typed), func(t *testing.T) {
			object := &mockObject{
				fn: func(encoder *gojay.Encoder) {
					AddTypeKey(encoder, "mock", typed, value)
				},
			}

			encoder := gojay.NewEncoder(ioutil.Discard)
			encoder.Encode(object)
		})
	}
}

func TestAddType(t *testing.T) {
	tests := map[types.Type]interface{}{
		types.Double:   float64(10),
		types.Int64:    int64(10),
		types.Uint64:   uint64(10),
		types.Fixed64:  uint64(10),
		types.Int32:    int32(10),
		types.Uint32:   uint32(10),
		types.Fixed32:  uint64(10),
		types.Float:    float32(10),
		types.String:   string("msg"),
		types.Bool:     true,
		types.Bytes:    []byte{10, 10},
		types.Sfixed32: int32(10),
		types.Sfixed64: int64(10),
		types.Sint32:   int32(10),
		types.Sint64:   int64(10),
	}

	for typed, value := range tests {
		t.Run(string(typed), func(t *testing.T) {
			object := &mockObject{
				fn: func(encoder *gojay.Encoder) {
					AddType(encoder, typed, value)
				},
			}

			encoder := gojay.NewEncoder(ioutil.Discard)
			encoder.Encode(object)
		})
	}
}

func TestDecodeTypeKey(t *testing.T) {
	tests := map[types.Type]interface{}{
		types.Double:   float64(10),
		types.Int64:    int64(10),
		types.Uint64:   uint64(10),
		types.Fixed64:  uint64(10),
		types.Int32:    int32(10),
		types.Uint32:   uint32(10),
		types.Fixed32:  uint64(10),
		types.Float:    float32(10),
		types.String:   string("msg"),
		types.Bool:     true,
		types.Bytes:    []byte{10, 10},
		types.Sfixed32: int32(10),
		types.Sfixed64: int64(10),
		types.Sint32:   int32(10),
		types.Sint64:   int64(10),
	}

	for typed, value := range tests {
		t.Run(string(typed), func(t *testing.T) {
			object := &mockObject{
				fn: func(encoder *gojay.Encoder) {
					AddTypeKey(encoder, "mock", typed, value)
				},
			}

			encoder := gojay.NewEncoder(ioutil.Discard)
			encoder.Encode(object)
		})
	}
}

func TestDecodeType(t *testing.T) {
	tests := map[types.Type]interface{}{
		types.Double:   float64(10),
		types.Int64:    int64(10),
		types.Uint64:   uint64(10),
		types.Fixed64:  uint64(10),
		types.Int32:    int32(10),
		types.Uint32:   uint32(10),
		types.Fixed32:  uint64(10),
		types.Float:    float32(10),
		types.String:   string("msg"),
		types.Bool:     true,
		types.Bytes:    []byte{10, 10},
		types.Sfixed32: int32(10),
		types.Sfixed64: int64(10),
		types.Sint32:   int32(10),
		types.Sint64:   int64(10),
	}

	for typed, value := range tests {
		t.Run(string(typed), func(t *testing.T) {
			object := &mockObject{
				fn: func(encoder *gojay.Encoder) {
					AddType(encoder, typed, value)
				},
			}

			encoder := gojay.NewEncoder(ioutil.Discard)
			encoder.Encode(object)
		})
	}
}
