package protobuffers

import (
	protobuf "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/template"
	"github.com/jhump/protoreflect/desc"
)

// NewSchema constructs a new schema manifest from the given file descriptors
func NewSchema(descriptors []*desc.FileDescriptor) specs.Schemas {
	result := make(specs.Schemas, 0)

	for _, descriptor := range descriptors {
		for _, message := range descriptor.GetMessageTypes() {
			result[message.GetFullyQualifiedName()] = NewMessage("", message)
		}
	}

	return result
}

// NewMessage constructs a schema Property with the given message descriptor
func NewMessage(path string, descriptor *desc.MessageDescriptor) *specs.Property {
	fields := descriptor.GetFields()
	result := &specs.Property{
		Path:        path,
		Name:        descriptor.GetFullyQualifiedName(),
		Description: descriptor.GetSourceInfo().GetLeadingComments(),
		Position:    1,
		Label:       labels.Optional,
		Template: specs.Template{
			Message: make(specs.Message, len(fields)),
		},
		Options: specs.Options{},
	}

	for _, field := range fields {
		result.Message[field.GetName()] = NewProperty(template.JoinPath(path, field.GetName()), field)
	}

	return result
}

// NewProperty constructs a schema Property with the given field descriptor
func NewProperty(path string, descriptor *desc.FieldDescriptor) *specs.Property {
	result := &specs.Property{
		Path:        path,
		Name:        descriptor.GetName(),
		Description: descriptor.GetSourceInfo().GetLeadingComments(),
		Position:    descriptor.GetNumber(),
		Options:     specs.Options{},
		Label:       Labels[descriptor.GetLabel()],
	}

	switch {
	case descriptor.GetType() == protobuf.FieldDescriptorProto_TYPE_ENUM:
		enum := descriptor.GetEnumType()
		keys := map[string]*specs.EnumValue{}
		positions := map[int32]*specs.EnumValue{}

		for _, value := range enum.GetValues() {
			result := &specs.EnumValue{
				Key:         value.GetName(),
				Position:    value.GetNumber(),
				Description: value.GetSourceInfo().GetLeadingComments(),
			}

			keys[value.GetName()] = result
			positions[value.GetNumber()] = result
		}

		result.Enum = &specs.Enum{
			Name:        enum.GetName(),
			Description: enum.GetSourceInfo().GetLeadingComments(),
			Keys:        keys,
			Positions:   positions,
		}

		break
	case descriptor.GetType() == protobuf.FieldDescriptorProto_TYPE_MESSAGE:
		var fields = descriptor.GetMessageType().GetFields()
		result.Message = make(specs.Message, len(fields))

		for _, field := range fields {
			result.Message[field.GetName()] = NewProperty(template.JoinPath(path, field.GetName()), field)
		}

		break
	default:
		result.Scalar = &specs.Scalar{
			Type: Types[descriptor.GetType()],
		}
	}

	if descriptor.GetLabel() == protobuf.FieldDescriptorProto_LABEL_REPEATED {
		result.Label = labels.Optional
		result.Template = specs.Template{
			Repeated: specs.Repeated{result.Template},
		}
	}

	return result
}
