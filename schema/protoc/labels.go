package protoc

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jexia/maestro/specs/types"
)

// ProtoLabels is a lookup table for field label types
var ProtoLabels = map[types.Label]descriptor.FieldDescriptorProto_Label{
	types.LabelOptional: descriptor.FieldDescriptorProto_LABEL_OPTIONAL,
	types.LabelRequired: descriptor.FieldDescriptorProto_LABEL_REQUIRED,
	types.LabelRepeated: descriptor.FieldDescriptorProto_LABEL_REPEATED,
}

// Labels is a lookup table for field label types
var Labels = map[descriptor.FieldDescriptorProto_Label]types.Label{
	descriptor.FieldDescriptorProto_LABEL_OPTIONAL: types.LabelOptional,
	descriptor.FieldDescriptorProto_LABEL_REQUIRED: types.LabelRequired,
	descriptor.FieldDescriptorProto_LABEL_REPEATED: types.LabelRepeated,
}
