package protobuffers

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jexia/semaphore/pkg/specs/labels"
)

const (
	// Optional label representing a optional field
	Optional = descriptor.FieldDescriptorProto_LABEL_OPTIONAL
	// Required label representing a required field
	Required = descriptor.FieldDescriptorProto_LABEL_REQUIRED
	// Repeated label representing a repeated field
	Repeated = descriptor.FieldDescriptorProto_LABEL_REPEATED
)

// ProtoLabels is a lookup table for field label types
var ProtoLabels = map[labels.Label]descriptor.FieldDescriptorProto_Label{
	labels.Optional: Optional,
	labels.Required: Required,
}

// Labels is a lookup table for field label types
var Labels = map[descriptor.FieldDescriptorProto_Label]labels.Label{
	Optional: labels.Optional,
	Required: labels.Required,
}
