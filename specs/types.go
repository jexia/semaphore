package specs

import (
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

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
	TypeGroup    Type = "group"
	TypeMessage  Type = "message"
	TypeBytes    Type = "bytes"
	TypeUint32   Type = "uint32"
	TypeEnum     Type = "enum"
	TypeSfixed32 Type = "sfixed32"
	TypeSfixed64 Type = "sfixed64"
	TypeSint32   Type = "sint32"
	TypeSint64   Type = "sint64"
)

// ProtoTypes is a lookup table for field descriptor types
var ProtoTypes = map[Type]descriptor.FieldDescriptorProto_Type{
	TypeDouble:   descriptor.FieldDescriptorProto_TYPE_DOUBLE,
	TypeFloat:    descriptor.FieldDescriptorProto_TYPE_FLOAT,
	TypeInt64:    descriptor.FieldDescriptorProto_TYPE_INT64,
	TypeUint64:   descriptor.FieldDescriptorProto_TYPE_UINT64,
	TypeInt32:    descriptor.FieldDescriptorProto_TYPE_INT32,
	TypeFixed64:  descriptor.FieldDescriptorProto_TYPE_FIXED64,
	TypeFixed32:  descriptor.FieldDescriptorProto_TYPE_FIXED32,
	TypeBool:     descriptor.FieldDescriptorProto_TYPE_BOOL,
	TypeString:   descriptor.FieldDescriptorProto_TYPE_STRING,
	TypeGroup:    descriptor.FieldDescriptorProto_TYPE_GROUP,
	TypeMessage:  descriptor.FieldDescriptorProto_TYPE_MESSAGE,
	TypeBytes:    descriptor.FieldDescriptorProto_TYPE_BYTES,
	TypeUint32:   descriptor.FieldDescriptorProto_TYPE_UINT32,
	TypeEnum:     descriptor.FieldDescriptorProto_TYPE_ENUM,
	TypeSfixed32: descriptor.FieldDescriptorProto_TYPE_SFIXED32,
	TypeSfixed64: descriptor.FieldDescriptorProto_TYPE_SFIXED64,
	TypeSint32:   descriptor.FieldDescriptorProto_TYPE_SINT32,
	TypeSint64:   descriptor.FieldDescriptorProto_TYPE_SINT64,
}

// GetType returns the type of the given proto type
func GetType(t descriptor.FieldDescriptorProto_Type) Type {
	for key, pt := range ProtoTypes {
		if pt == t {
			return key
		}
	}

	return "unkown"
}

// IsType checks whether the given value is a type definition
func IsType(value string) bool {
	return strings.HasPrefix(value, TypeOpen) && strings.HasSuffix(value, TypeClose)
}

// GetTypeContent trims the opening and closing tags from the given type value
func GetTypeContent(value string) string {
	value = strings.Replace(value, TypeOpen, "", 1)
	value = strings.Replace(value, TypeClose, "", 1)
	value = strings.TrimSpace(value)
	return value
}

// SetType parses the given type and sets the property type
func SetType(property *Property, value cty.Value) {
	if value.Type() != cty.String {
		return
	}

	property.Type = Type(GetTypeContent(value.AsString()))
}

// SetDefaultValue sets the given value as default value inside the given property
func SetDefaultValue(property *Property, value cty.Value) {
	switch value.Type() {
	case cty.String:
		property.Default = value.AsString()
		property.Type = TypeString
	case cty.Number:
		var def int64
		gocty.FromCtyValue(value, &def)

		property.Default = def
		property.Type = TypeInt64
	case cty.Bool:
		var def bool
		gocty.FromCtyValue(value, &def)

		property.Default = def
		property.Type = TypeBool
	}
}
