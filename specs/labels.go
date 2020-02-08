package specs

import "github.com/golang/protobuf/protoc-gen-go/descriptor"

// Label represents a value label
type Label string

// Spec labels
const (
	LabelOptional Label = "optional"
	LabelRequired Label = "required"
	LabelRepeated Label = "repeated"
)

// ProtoLabels is a lookup table for field label types
var ProtoLabels = map[Label]descriptor.FieldDescriptorProto_Label{
	LabelOptional: descriptor.FieldDescriptorProto_LABEL_OPTIONAL,
	LabelRequired: descriptor.FieldDescriptorProto_LABEL_REQUIRED,
	LabelRepeated: descriptor.FieldDescriptorProto_LABEL_REPEATED,
}
