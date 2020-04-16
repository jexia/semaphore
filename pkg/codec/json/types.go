package json

import (
	"encoding/base64"

	"github.com/francoispqt/gojay"
	"github.com/jexia/maestro/pkg/specs/types"
)

// AddTypeKey encodes the given value into the given encoder
func AddTypeKey(encoder *gojay.Encoder, key string, typed types.Type, value interface{}) {
	switch typed {
	case types.Double:
		encoder.AddFloat64Key(key, Float64Empty(value))
	case types.Int64:
		encoder.AddInt64Key(key, Int64Empty(value))
	case types.Uint64:
		encoder.AddUint64Key(key, Uint64Empty(value))
	case types.Fixed64:
		encoder.AddUint64Key(key, Uint64Empty(value))
	case types.Int32:
		encoder.AddInt32Key(key, Int32Empty(value))
	case types.Uint32:
		encoder.AddUint32Key(key, Uint32Empty(value))
	case types.Fixed32:
		encoder.AddUint64Key(key, Uint64Empty(value))
	case types.Float:
		encoder.AddFloat32Key(key, Float32Empty(value))
	case types.String:
		encoder.AddStringKey(key, StringEmpty(value))
	case types.Bool:
		encoder.AddBoolKey(key, BoolEmpty(value))
	case types.Bytes:
		encoder.AddStringKey(key, BytesBase64Empty(value))
	case types.Sfixed32:
		encoder.AddInt32Key(key, Int32Empty(value))
	case types.Sfixed64:
		encoder.AddInt64Key(key, Int64Empty(value))
	case types.Sint32:
		encoder.AddInt32Key(key, Int32Empty(value))
	case types.Sint64:
		encoder.AddInt64Key(key, Int64Empty(value))
	}
}

// AddType encodes the given value into the given encoder
func AddType(encoder *gojay.Encoder, typed types.Type, value interface{}) {
	switch typed {
	case types.Double:
		encoder.AddFloat64(Float64Empty(value))
	case types.Int64:
		encoder.AddInt64(Int64Empty(value))
	case types.Uint64:
		encoder.AddUint64(Uint64Empty(value))
	case types.Fixed64:
		encoder.AddUint64(Uint64Empty(value))
	case types.Int32:
		encoder.AddInt32(Int32Empty(value))
	case types.Uint32:
		encoder.AddUint32(Uint32Empty(value))
	case types.Fixed32:
		encoder.AddUint64(Uint64Empty(value))
	case types.Float:
		encoder.AddFloat32(Float32Empty(value))
	case types.String:
		encoder.AddString(StringEmpty(value))
	case types.Bool:
		encoder.AddBool(BoolEmpty(value))
	case types.Bytes:
		encoder.AddString(BytesBase64Empty(value))
	case types.Sfixed32:
		encoder.AddInt32(Int32Empty(value))
	case types.Sfixed64:
		encoder.AddInt64(Int64Empty(value))
	case types.Sint32:
		encoder.AddInt32(Int32Empty(value))
	case types.Sint64:
		encoder.AddInt64(Int64Empty(value))
	}
}

// DecodeType decodes the given property from the given decoder
func DecodeType(decoder *gojay.Decoder, prop types.Type) interface{} {
	switch prop {
	case types.Double:
		var value float64
		decoder.AddFloat64(&value)
		return value
	case types.Float:
		var value float32
		decoder.AddFloat32(&value)
		return value
	case types.Int64:
		var value int64
		decoder.AddInt64(&value)
		return value
	case types.Uint64:
		var value uint64
		decoder.AddUint64(&value)
		return value
	case types.Fixed64:
		var value uint64
		decoder.AddUint64(&value)
		return value
	case types.Int32:
		var value int32
		decoder.AddInt32(&value)
		return value
	case types.Uint32:
		var value uint32
		decoder.AddUint32(&value)
		return value
	case types.Fixed32:
		var value uint64
		decoder.AddUint64(&value)
		return value
	case types.String:
		var value string
		decoder.AddString(&value)
		return value
	case types.Bool:
		var value bool
		decoder.AddBool(&value)
		return value
	case types.Bytes:
		var raw string
		decoder.AddString(&raw)

		value := make([]byte, len(raw))
		base64.StdEncoding.Decode(value, []byte(raw))
		return value
	case types.Sfixed32:
		var value int32
		decoder.AddInt32(&value)
		return value
	case types.Sfixed64:
		var value int64
		decoder.AddInt64(&value)
		return value
	case types.Sint32:
		var value int32
		decoder.AddInt32(&value)
		return value
	case types.Sint64:
		var value int64
		decoder.AddInt64(&value)
		return value
	}

	return nil
}

// StringEmpty returns the given value as a string or a empty string if the value is nil
func StringEmpty(val interface{}) string {
	if val == nil {
		return ""
	}

	return val.(string)
}

// BoolEmpty returns the given value as a bool or a empty bool if the value is nil
func BoolEmpty(val interface{}) bool {
	if val == nil {
		return false
	}

	return val.(bool)
}

// Int32Empty returns the given value as a int32 or a empty int32 if the value is nil
func Int32Empty(val interface{}) int32 {
	if val == nil {
		return 0
	}

	return val.(int32)
}

// Uint32Empty returns the given value as a uint32 or a empty uint32 if the value is nil
func Uint32Empty(val interface{}) uint32 {
	if val == nil {
		return 0
	}

	return val.(uint32)
}

// Int64Empty returns the given value as a int64 or a empty int64 if the value is nil
func Int64Empty(val interface{}) int64 {
	if val == nil {
		return 0
	}

	return val.(int64)
}

// Uint64Empty returns the given value as a uint64 or a empty uint64 if the value is nil
func Uint64Empty(val interface{}) uint64 {
	if val == nil {
		return 0
	}

	return val.(uint64)
}

// Float64Empty returns the given value as a float64 or a empty float64 if the value is nil
func Float64Empty(val interface{}) float64 {
	if val == nil {
		return 0
	}

	return val.(float64)
}

// Float32Empty returns the given value as a float32 or a empty float32 if the value is nil
func Float32Empty(val interface{}) float32 {
	if val == nil {
		return 0
	}

	return val.(float32)
}

// BytesBase64Empty returns the given bytes buffer as a base64 string or a empty string if the value is nil
func BytesBase64Empty(val interface{}) string {
	if val == nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(val.([]byte))
}
