package protobuffers

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jexia/semaphore/v2/pkg/specs/types"
)

// ProtoTypes is a lookup table for field descriptor types
var ProtoTypes = map[types.Type]descriptor.FieldDescriptorProto_Type{
	types.Double:   descriptor.FieldDescriptorProto_TYPE_DOUBLE,
	types.Float:    descriptor.FieldDescriptorProto_TYPE_FLOAT,
	types.Int64:    descriptor.FieldDescriptorProto_TYPE_INT64,
	types.Uint64:   descriptor.FieldDescriptorProto_TYPE_UINT64,
	types.Int32:    descriptor.FieldDescriptorProto_TYPE_INT32,
	types.Fixed64:  descriptor.FieldDescriptorProto_TYPE_FIXED64,
	types.Fixed32:  descriptor.FieldDescriptorProto_TYPE_FIXED32,
	types.Bool:     descriptor.FieldDescriptorProto_TYPE_BOOL,
	types.String:   descriptor.FieldDescriptorProto_TYPE_STRING,
	types.Message:  descriptor.FieldDescriptorProto_TYPE_MESSAGE,
	types.Bytes:    descriptor.FieldDescriptorProto_TYPE_BYTES,
	types.Uint32:   descriptor.FieldDescriptorProto_TYPE_UINT32,
	types.Enum:     descriptor.FieldDescriptorProto_TYPE_ENUM,
	types.Sfixed32: descriptor.FieldDescriptorProto_TYPE_SFIXED32,
	types.Sfixed64: descriptor.FieldDescriptorProto_TYPE_SFIXED64,
	types.Sint32:   descriptor.FieldDescriptorProto_TYPE_SINT32,
	types.Sint64:   descriptor.FieldDescriptorProto_TYPE_SINT64,
}

// Types is a lookup table for field descriptor types
var Types = map[descriptor.FieldDescriptorProto_Type]types.Type{
	descriptor.FieldDescriptorProto_TYPE_DOUBLE:   types.Double,
	descriptor.FieldDescriptorProto_TYPE_FLOAT:    types.Float,
	descriptor.FieldDescriptorProto_TYPE_INT64:    types.Int64,
	descriptor.FieldDescriptorProto_TYPE_UINT64:   types.Uint64,
	descriptor.FieldDescriptorProto_TYPE_INT32:    types.Int32,
	descriptor.FieldDescriptorProto_TYPE_FIXED64:  types.Fixed64,
	descriptor.FieldDescriptorProto_TYPE_FIXED32:  types.Fixed32,
	descriptor.FieldDescriptorProto_TYPE_BOOL:     types.Bool,
	descriptor.FieldDescriptorProto_TYPE_STRING:   types.String,
	descriptor.FieldDescriptorProto_TYPE_MESSAGE:  types.Message,
	descriptor.FieldDescriptorProto_TYPE_BYTES:    types.Bytes,
	descriptor.FieldDescriptorProto_TYPE_UINT32:   types.Uint32,
	descriptor.FieldDescriptorProto_TYPE_ENUM:     types.Enum,
	descriptor.FieldDescriptorProto_TYPE_SFIXED32: types.Sfixed32,
	descriptor.FieldDescriptorProto_TYPE_SFIXED64: types.Sfixed64,
	descriptor.FieldDescriptorProto_TYPE_SINT32:   types.Sint32,
	descriptor.FieldDescriptorProto_TYPE_SINT64:   types.Sint64,
}
