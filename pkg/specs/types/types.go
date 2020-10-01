package types

import (
	"encoding/base64"
	"strconv"
)

// Type represents a value type definition.
type Type string

// Spec types.
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
	Array    Type = "array"
	Message  Type = "message"
	Bytes    Type = "bytes"
	Uint32   Type = "uint32"
	Enum     Type = "enum"
	Sfixed32 Type = "sfixed32"
	Sfixed64 Type = "sfixed64"
	Sint32   Type = "sint32"
	Sint64   Type = "sint64"
	Unknown  Type = "unknown"
)

// DecodeFromString decodes the given property from string.
func DecodeFromString(raw string, typed Type) (interface{}, error) {
	switch typed {
	case Double:
		return strconv.ParseFloat(raw, 64)
	case Float:
		value, err := strconv.ParseFloat(raw, 64)
		return float64(value), err
	case Int64:
		return strconv.ParseInt(raw, 10, 64)
	case Uint64:
		return strconv.ParseUint(raw, 10, 64)
	case Fixed64:
		return strconv.ParseUint(raw, 10, 64)
	case Int32:
		value, err := strconv.ParseInt(raw, 10, 32)
		return int32(value), err
	case Uint32:
		value, err := strconv.ParseUint(raw, 10, 32)
		return uint32(value), err
	case Fixed32:
		value, err := strconv.ParseUint(raw, 10, 32)
		return uint32(value), err
	case String:
		return raw, nil
	case Bool:
		return strconv.ParseBool(raw)
	case Bytes:
		return base64.StdEncoding.DecodeString(raw)
	case Sfixed32:
		value, err := strconv.ParseInt(raw, 10, 32)
		return int32(value), err
	case Sfixed64:
		value, err := strconv.ParseInt(raw, 10, 64)
		return value, err
	case Sint32:
		value, err := strconv.ParseInt(raw, 10, 32)
		return int32(value), err
	case Sint64:
		return strconv.ParseInt(raw, 10, 64)
	default:
		return nil, ErrUnknownType(typed)
	}
}
