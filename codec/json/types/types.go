package types

import (
	"github.com/francoispqt/gojay"
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/types"
)

func Add(encoder *gojay.Encoder, key string, typed types.Type, value interface{}) {
	switch typed {
	case types.TypeDouble:
		encoder.AddFloat64Key(key, Float64Empty(value))
	case types.TypeInt32:
		encoder.AddInt32Key(key, Int32Empty(value))
	case types.TypeInt64:
		encoder.AddInt64Key(key, Int64Empty(value))
	case types.TypeFloat:
		encoder.AddFloat32Key(key, Float32Empty(value))
	case types.TypeString:
		encoder.AddStringKey(key, StringEmpty(value))
	case types.TypeBool:
		encoder.AddBoolKey(key, BoolEmpty(value))
	}
}

func Decode(decoder *gojay.Decoder, prop *specs.Property, store *refs.Store) interface{} {
	switch prop.GetType() {
	case types.TypeDouble:
		var value float64
		decoder.AddFloat64(&value)
		return value
	case types.TypeFloat:
		var value float32
		decoder.AddFloat32(&value)
		return value
	case types.TypeInt32:
		var value int32
		decoder.AddInt32(&value)
		return value
	case types.TypeInt64:
		var value int64
		decoder.AddInt64(&value)
		return value
	case types.TypeString:
		var value string
		decoder.AddString(&value)
		return value
	case types.TypeBool:
		var value bool
		decoder.AddBool(&value)
		return value
	}

	return nil
}

func StringEmpty(val interface{}) string {
	if val == nil {
		return ""
	}

	return val.(string)
}

func BoolEmpty(val interface{}) bool {
	if val == nil {
		return false
	}

	return val.(bool)
}

func Int32Empty(val interface{}) int32 {
	if val == nil {
		return 0
	}

	return val.(int32)
}

func Int64Empty(val interface{}) int64 {
	if val == nil {
		return 0
	}

	return val.(int64)
}

func Float64Empty(val interface{}) float64 {
	if val == nil {
		return 0
	}

	return val.(float64)
}

func Float32Empty(val interface{}) float32 {
	if val == nil {
		return 0
	}

	return val.(float32)
}
