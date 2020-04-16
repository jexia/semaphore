package types

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
