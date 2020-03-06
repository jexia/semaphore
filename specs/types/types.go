package types

const (
	// TypeOpen tag
	TypeOpen = "<"
	// TypeClose tag
	TypeClose = ">"
)

// Type represents a value type definition
type Type string

// Spec types
const (
	TypeDouble   Type = "double"
	TypeFloat    Type = "float"
	TypeInt64    Type = "int64"
	TypeUint64   Type = "uint64"
	TypeInt32    Type = "int32"
	TypeFixed64  Type = "fixed64"
	TypeFixed32  Type = "fixed32"
	TypeBool     Type = "bool"
	TypeString   Type = "string"
	TypeMessage  Type = "message"
	TypeBytes    Type = "bytes"
	TypeUint32   Type = "uint32"
	TypeEnum     Type = "enum"
	TypeSfixed32 Type = "sfixed32"
	TypeSfixed64 Type = "sfixed64"
	TypeSint32   Type = "sint32"
	TypeSint64   Type = "sint64"
)
