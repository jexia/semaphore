package json

import (
	"encoding/base64"

	"github.com/francoispqt/gojay"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/types"
)

// AddType encodes the given value into the given encoder
func AddType(encoder *gojay.Encoder, key string, typed types.Type, value interface{}) {
	switch typed {
	case types.TypeDouble:
		encoder.AddFloat64Key(key, Float64Empty(value))
	case types.TypeInt64:
		encoder.AddInt64Key(key, Int64Empty(value))
	case types.TypeUint64:
		encoder.AddUint64Key(key, Uint64Empty(value))
	case types.TypeFixed64:
		encoder.AddUint64Key(key, Uint64Empty(value))
	case types.TypeInt32:
		encoder.AddInt32Key(key, Int32Empty(value))
	case types.TypeUint32:
		encoder.AddUint32Key(key, Uint32Empty(value))
	case types.TypeFixed32:
		encoder.AddUint64Key(key, Uint64Empty(value))
	case types.TypeFloat:
		encoder.AddFloat32Key(key, Float32Empty(value))
	case types.TypeString:
		encoder.AddStringKey(key, StringEmpty(value))
	case types.TypeBool:
		encoder.AddBoolKey(key, BoolEmpty(value))
	case types.TypeBytes:
		encoder.AddStringKey(key, BytesBase64Empty(value))
	case types.TypeSfixed32:
		encoder.AddInt32Key(key, Int32Empty(value))
	case types.TypeSfixed64:
		encoder.AddInt64Key(key, Int64Empty(value))
	case types.TypeSint32:
		encoder.AddInt32Key(key, Int32Empty(value))
	case types.TypeSint64:
		encoder.AddInt64Key(key, Int64Empty(value))
	}
}

// DecodeType decodes the given property from the given decoder
func DecodeType(decoder *gojay.Decoder, prop *specs.Property) interface{} {
	switch prop.Type {
	case types.TypeDouble:
		var value float64
		decoder.AddFloat64(&value)
		return value
	case types.TypeFloat:
		var value float32
		decoder.AddFloat32(&value)
		return value
	case types.TypeInt64:
		var value int64
		decoder.AddInt64(&value)
		return value
	case types.TypeUint64:
		var value uint64
		decoder.AddUint64(&value)
		return value
	case types.TypeFixed64:
		var value uint64
		decoder.AddUint64(&value)
		return value
	case types.TypeInt32:
		var value int32
		decoder.AddInt32(&value)
		return value
	case types.TypeUint32:
		var value uint32
		decoder.AddUint32(&value)
		return value
	case types.TypeFixed32:
		var value uint64
		decoder.AddUint64(&value)
		return value
	case types.TypeString:
		var value string
		decoder.AddString(&value)
		return value
	case types.TypeBool:
		var value bool
		decoder.AddBool(&value)
		return value
	case types.TypeBytes:
		var raw string
		decoder.AddString(&raw)

		value := make([]byte, len(raw))
		base64.StdEncoding.Decode(value, []byte(raw))
		return value
	case types.TypeSfixed32:
		var value int32
		decoder.AddInt32(&value)
		return value
	case types.TypeSfixed64:
		var value int64
		decoder.AddInt64(&value)
		return value
	case types.TypeSint32:
		var value int32
		decoder.AddInt32(&value)
		return value
	case types.TypeSint64:
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
