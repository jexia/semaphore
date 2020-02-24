package types

import (
	"github.com/francoispqt/gojay"
	"github.com/jexia/maestro/specs/types"
)

func Add(encoder *gojay.Encoder, key string, typed types.Type, value interface{}) {
	switch typed {
	case types.TypeDouble:
		encoder.AddFloat64Key(key, Float64Empty(value))
	case types.TypeFloat:
		encoder.AddFloat32Key(key, Float32Empty(value))
	case types.TypeString:
		encoder.AddStringKey(key, StringEmpty(value))
	case types.TypeBool:
		encoder.AddBoolKey(key, BoolEmpty(value))
	}
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