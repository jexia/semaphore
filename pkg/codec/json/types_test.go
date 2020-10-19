package json

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/francoispqt/gojay"
	"github.com/jexia/semaphore/pkg/specs/types"
)

type mockObject struct {
	enc func(*gojay.Encoder)
	dec func(*gojay.Decoder)
}

func (m *mockObject) UnmarshalJSONObject(decoder *gojay.Decoder, _ string) error {
	m.dec(decoder)
	return nil
}

func (m *mockObject) NKeys() int {
	return 1
}

func (m *mockObject) MarshalJSONObject(encoder *gojay.Encoder) {
	m.enc(encoder)
}

func (m *mockObject) IsNil() bool {
	return m == nil
}

func TestAddTypeKey(t *testing.T) {
	type test struct {
		dataType types.Type
		value    interface{}
	}

	tests := map[string]test{
		fmt.Sprintf("%s from float32", string(types.Double)): {types.Double, float32(10)},
		fmt.Sprintf("%s from float64", string(types.Double)): {types.Double, float64(10)},
		string(types.Int64):   {types.Int64, int64(10)},
		string(types.Uint64):  {types.Uint64, uint64(10)},
		string(types.Fixed64): {types.Fixed64, uint64(10)},
		string(types.Int32):   {types.Int32, int32(10)},
		string(types.Uint32):  {types.Uint32, uint32(10)},
		string(types.Fixed32): {types.Fixed32, uint64(10)},
		fmt.Sprintf("%s from float32", string(types.Float)): {types.Float, float32(10)},
		fmt.Sprintf("%s from float64", string(types.Float)): {types.Float, float64(10)},
		string(types.Float):    {types.Float, float32(10)},
		string(types.String):   {types.String, string("msg")},
		string(types.Bool):     {types.Bool, true},
		string(types.Bytes):    {types.Bytes, []byte{10, 10}},
		string(types.Sfixed32): {types.Sfixed32, int32(10)},
		string(types.Sfixed64): {types.Sfixed64, int64(10)},
		string(types.Sint32):   {types.Sint32, int32(10)},
		string(types.Sint64):   {types.Sint64, int64(10)},
	}

	for title, test := range tests {
		t.Run(title, func(t *testing.T) {
			object := &mockObject{
				enc: func(encoder *gojay.Encoder) {
					AddTypeKey(encoder, "mock", test.dataType, test.value)
				},
			}

			encoder := gojay.NewEncoder(ioutil.Discard)
			err := encoder.Encode(object)
			if err != nil {
				t.Fatal(err)
			}
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
				enc: func(encoder *gojay.Encoder) {
					AddType(encoder, typed, value)
				},
			}

			encoder := gojay.NewEncoder(ioutil.Discard)
			err := encoder.Encode(object)
			if err != nil {
				t.Fatal(err)
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
			object := &mockObject{
				enc: func(encoder *gojay.Encoder) {
					AddType(encoder, typed, value)
				},
			}

			encoder := gojay.NewEncoder(ioutil.Discard)
			err := encoder.Encode(object)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestDecodeType(t *testing.T) {
	tests := map[types.Type]string{
		types.Double:   `{"mock":10}`,
		types.Int64:    `{"mock":10}`,
		types.Uint64:   `{"mock":10}`,
		types.Fixed64:  `{"mock":10}`,
		types.Int32:    `{"mock":10}`,
		types.Uint32:   `{"mock":10}`,
		types.Fixed32:  `{"mock":10}`,
		types.Float:    `{"mock":10}`,
		types.String:   `{"mock":"msg"}`,
		types.Bool:     `{"mock":true}`,
		types.Bytes:    `{"mock":"aGVsbG8="}`,
		types.Sfixed32: `{"mock":10}`,
		types.Sfixed64: `{"mock":10}`,
		types.Sint32:   `{"mock":10}`,
		types.Sint64:   `{"mock":10}`,
	}

	for typed, value := range tests {
		t.Run(string(typed), func(t *testing.T) {
			object := &mockObject{
				dec: func(decoder *gojay.Decoder) {
					value, err := DecodeType(decoder, typed)
					if err != nil {
						t.Fatal(err)
					}

					if value == nil {
						t.Fatal("empty value returned, expected decoded value")
					}
				},
			}

			encoder := gojay.NewDecoder(strings.NewReader(value))
			err := encoder.Decode(object)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
