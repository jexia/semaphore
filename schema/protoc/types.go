package protoc

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jexia/maestro/specs"
)

// ProtoTypes is a lookup table for field descriptor types
var ProtoTypes = map[specs.Type]descriptor.FieldDescriptorProto_Type{
	specs.TypeDouble:   descriptor.FieldDescriptorProto_TYPE_DOUBLE,
	specs.TypeFloat:    descriptor.FieldDescriptorProto_TYPE_FLOAT,
	specs.TypeInt64:    descriptor.FieldDescriptorProto_TYPE_INT64,
	specs.TypeUint64:   descriptor.FieldDescriptorProto_TYPE_UINT64,
	specs.TypeInt32:    descriptor.FieldDescriptorProto_TYPE_INT32,
	specs.TypeFixed64:  descriptor.FieldDescriptorProto_TYPE_FIXED64,
	specs.TypeFixed32:  descriptor.FieldDescriptorProto_TYPE_FIXED32,
	specs.TypeBool:     descriptor.FieldDescriptorProto_TYPE_BOOL,
	specs.TypeString:   descriptor.FieldDescriptorProto_TYPE_STRING,
	specs.TypeGroup:    descriptor.FieldDescriptorProto_TYPE_GROUP,
	specs.TypeMessage:  descriptor.FieldDescriptorProto_TYPE_MESSAGE,
	specs.TypeBytes:    descriptor.FieldDescriptorProto_TYPE_BYTES,
	specs.TypeUint32:   descriptor.FieldDescriptorProto_TYPE_UINT32,
	specs.TypeEnum:     descriptor.FieldDescriptorProto_TYPE_ENUM,
	specs.TypeSfixed32: descriptor.FieldDescriptorProto_TYPE_SFIXED32,
	specs.TypeSfixed64: descriptor.FieldDescriptorProto_TYPE_SFIXED64,
	specs.TypeSint32:   descriptor.FieldDescriptorProto_TYPE_SINT32,
	specs.TypeSint64:   descriptor.FieldDescriptorProto_TYPE_SINT64,
}

// Types is a lookup table for field descriptor types
var Types = map[descriptor.FieldDescriptorProto_Type]specs.Type{
	descriptor.FieldDescriptorProto_TYPE_DOUBLE:   specs.TypeDouble,
	descriptor.FieldDescriptorProto_TYPE_FLOAT:    specs.TypeFloat,
	descriptor.FieldDescriptorProto_TYPE_INT64:    specs.TypeInt64,
	descriptor.FieldDescriptorProto_TYPE_UINT64:   specs.TypeUint64,
	descriptor.FieldDescriptorProto_TYPE_INT32:    specs.TypeInt32,
	descriptor.FieldDescriptorProto_TYPE_FIXED64:  specs.TypeFixed64,
	descriptor.FieldDescriptorProto_TYPE_FIXED32:  specs.TypeFixed32,
	descriptor.FieldDescriptorProto_TYPE_BOOL:     specs.TypeBool,
	descriptor.FieldDescriptorProto_TYPE_STRING:   specs.TypeString,
	descriptor.FieldDescriptorProto_TYPE_GROUP:    specs.TypeGroup,
	descriptor.FieldDescriptorProto_TYPE_MESSAGE:  specs.TypeMessage,
	descriptor.FieldDescriptorProto_TYPE_BYTES:    specs.TypeBytes,
	descriptor.FieldDescriptorProto_TYPE_UINT32:   specs.TypeUint32,
	descriptor.FieldDescriptorProto_TYPE_ENUM:     specs.TypeEnum,
	descriptor.FieldDescriptorProto_TYPE_SFIXED32: specs.TypeSfixed32,
	descriptor.FieldDescriptorProto_TYPE_SFIXED64: specs.TypeSfixed64,
	descriptor.FieldDescriptorProto_TYPE_SINT32:   specs.TypeSint32,
	descriptor.FieldDescriptorProto_TYPE_SINT64:   specs.TypeSint64,
}
