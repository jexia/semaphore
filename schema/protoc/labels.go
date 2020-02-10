package protoc

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jexia/maestro/specs"
)

// ProtoLabels is a lookup table for field label types
var ProtoLabels = map[specs.Label]descriptor.FieldDescriptorProto_Label{
	specs.LabelOptional: descriptor.FieldDescriptorProto_LABEL_OPTIONAL,
	specs.LabelRequired: descriptor.FieldDescriptorProto_LABEL_REQUIRED,
	specs.LabelRepeated: descriptor.FieldDescriptorProto_LABEL_REPEATED,
}

// Labels is a lookup table for field label types
var Labels = map[descriptor.FieldDescriptorProto_Label]specs.Label{
	descriptor.FieldDescriptorProto_LABEL_OPTIONAL: specs.LabelOptional,
	descriptor.FieldDescriptorProto_LABEL_REQUIRED: specs.LabelRequired,
	descriptor.FieldDescriptorProto_LABEL_REPEATED: specs.LabelRepeated,
}
