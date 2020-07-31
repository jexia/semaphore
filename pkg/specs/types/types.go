package types

import (
	"errors"
	"reflect"
)

// Type represents a value type definition
type Type string

// Spec types
const (
	Double   Type = "double"
	Float    Type = "float"
	Int64    Type = "int64"
	Uint64   Type = "uint64"
	Int32    Type = "int32"
	Fixed64  Type = "fixed64"
	Fixed32  Type = "fixed32"
	Bool     Type = "bool"
	String   Type = "string"
	Message  Type = "message"
	Bytes    Type = "bytes"
	Uint32   Type = "uint32"
	Enum     Type = "enum"
	Sfixed32 Type = "sfixed32"
	Sfixed64 Type = "sfixed64"
	Sint32   Type = "sint32"
	Sint64   Type = "sint64"
)

// Ensure ensures that the given data type stays consistent
func Ensure(t Type, v interface{}) (interface{}, error) {
	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)

	switch t {
	case Double:
		return convert(val, typ, reflect.TypeOf(float64(0)))
	case Float:
		return convert(val, typ, reflect.TypeOf(float32(0)))
	case Int64:
		return convert(val, typ, reflect.TypeOf(int64(0)))
	case Uint64:
		return convert(val, typ, reflect.TypeOf(uint64(0)))
	}

	return v, nil
}

func convert(val reflect.Value, org reflect.Type, typ reflect.Type) (interface{}, error) {
	if !org.ConvertibleTo(typ) {
		return nil, errors.New("type is not convertible to " + typ.Name())
	}

	return val.Convert(typ).Interface(), nil
}
