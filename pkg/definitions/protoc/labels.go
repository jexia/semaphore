package protoc

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jexia/maestro/pkg/specs/labels"
)

// ProtoLabels is a lookup table for field label types
var ProtoLabels = map[labels.Label]descriptor.FieldDescriptorProto_Label{
	labels.Optional: descriptor.FieldDescriptorProto_LABEL_OPTIONAL,
	labels.Required: descriptor.FieldDescriptorProto_LABEL_REQUIRED,
	labels.Repeated: descriptor.FieldDescriptorProto_LABEL_REPEATED,
}

// Labels is a lookup table for field label types
var Labels = map[descriptor.FieldDescriptorProto_Label]labels.Label{
	descriptor.FieldDescriptorProto_LABEL_OPTIONAL: labels.Optional,
	descriptor.FieldDescriptorProto_LABEL_REQUIRED: labels.Required,
	descriptor.FieldDescriptorProto_LABEL_REPEATED: labels.Repeated,
}
