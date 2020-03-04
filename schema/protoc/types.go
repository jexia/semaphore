package protoc

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jexia/maestro/specs/types"
)

// ProtoTypes is a lookup table for field descriptor types
var ProtoTypes = map[types.Type]descriptor.FieldDescriptorProto_Type{
	types.TypeDouble:   descriptor.FieldDescriptorProto_TYPE_DOUBLE,
	types.TypeFloat:    descriptor.FieldDescriptorProto_TYPE_FLOAT,
	types.TypeInt64:    descriptor.FieldDescriptorProto_TYPE_INT64,
	types.TypeUint64:   descriptor.FieldDescriptorProto_TYPE_UINT64,
	types.TypeInt32:    descriptor.FieldDescriptorProto_TYPE_INT32,
	types.TypeFixed64:  descriptor.FieldDescriptorProto_TYPE_FIXED64,
	types.TypeFixed32:  descriptor.FieldDescriptorProto_TYPE_FIXED32,
	types.TypeBool:     descriptor.FieldDescriptorProto_TYPE_BOOL,
	types.TypeString:   descriptor.FieldDescriptorProto_TYPE_STRING,
	types.TypeGroup:    descriptor.FieldDescriptorProto_TYPE_GROUP,
	types.TypeMessage:  descriptor.FieldDescriptorProto_TYPE_MESSAGE,
	types.TypeBytes:    descriptor.FieldDescriptorProto_TYPE_BYTES,
	types.TypeUint32:   descriptor.FieldDescriptorProto_TYPE_UINT32,
	types.TypeEnum:     descriptor.FieldDescriptorProto_TYPE_ENUM,
	types.TypeSfixed32: descriptor.FieldDescriptorProto_TYPE_SFIXED32,
	types.TypeSfixed64: descriptor.FieldDescriptorProto_TYPE_SFIXED64,
	types.TypeSint32:   descriptor.FieldDescriptorProto_TYPE_SINT32,
	types.TypeSint64:   descriptor.FieldDescriptorProto_TYPE_SINT64,
}

// Types is a lookup table for field descriptor types
var Types = map[descriptor.FieldDescriptorProto_Type]types.Type{
	descriptor.FieldDescriptorProto_TYPE_DOUBLE:   types.TypeDouble,
	descriptor.FieldDescriptorProto_TYPE_FLOAT:    types.TypeFloat,
	descriptor.FieldDescriptorProto_TYPE_INT64:    types.TypeInt64,
	descriptor.FieldDescriptorProto_TYPE_UINT64:   types.TypeUint64,
	descriptor.FieldDescriptorProto_TYPE_INT32:    types.TypeInt32,
	descriptor.FieldDescriptorProto_TYPE_FIXED64:  types.TypeFixed64,
	descriptor.FieldDescriptorProto_TYPE_FIXED32:  types.TypeFixed32,
	descriptor.FieldDescriptorProto_TYPE_BOOL:     types.TypeBool,
	descriptor.FieldDescriptorProto_TYPE_STRING:   types.TypeString,
	descriptor.FieldDescriptorProto_TYPE_GROUP:    types.TypeGroup,
	descriptor.FieldDescriptorProto_TYPE_MESSAGE:  types.TypeMessage,
	descriptor.FieldDescriptorProto_TYPE_BYTES:    types.TypeBytes,
	descriptor.FieldDescriptorProto_TYPE_UINT32:   types.TypeUint32,
	descriptor.FieldDescriptorProto_TYPE_ENUM:     types.TypeEnum,
	descriptor.FieldDescriptorProto_TYPE_SFIXED32: types.TypeSfixed32,
	descriptor.FieldDescriptorProto_TYPE_SFIXED64: types.TypeSfixed64,
	descriptor.FieldDescriptorProto_TYPE_SINT32:   types.TypeSint32,
	descriptor.FieldDescriptorProto_TYPE_SINT64:   types.TypeSint64,
}
