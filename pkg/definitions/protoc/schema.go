package protoc

import (
	protobuf "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/labels"
	"github.com/jexia/maestro/pkg/specs/template"
	"github.com/jexia/maestro/pkg/specs/types"
	"github.com/jhump/protoreflect/desc"
)

// NewSchema constructs a new schema manifest from the given file descriptors
func NewSchema(descriptors []*desc.FileDescriptor) *specs.SchemaManifest {
	result := &specs.SchemaManifest{
		Properties: make(map[string]*specs.Property, 0),
	}

	for _, descriptor := range descriptors {
		for _, message := range descriptor.GetMessageTypes() {
			result.Properties[message.GetFullyQualifiedName()] = NewMessage("", message)
		}
	}

	return result
}

// NewMessage constructs a schema Property with the given message descriptor
func NewMessage(path string, descriptor *desc.MessageDescriptor) *specs.Property {
	fields := descriptor.GetFields()
	result := &specs.Property{
		Path:     path,
		Name:     descriptor.GetFullyQualifiedName(),
		Comment:  descriptor.GetSourceInfo().GetLeadingComments(),
		Position: 1,
		Type:     types.Message,
		Label:    labels.Optional,
		Nested:   make(map[string]*specs.Property, len(fields)),
		Options:  specs.Options{},
	}

	for _, field := range fields {
		result.Nested[field.GetName()] = NewProperty(template.JoinPath(path, field.GetName()), field)
	}

	return result
}

// NewProperty constructs a schema Property with the given field descriptor
func NewProperty(path string, descriptor *desc.FieldDescriptor) *specs.Property {
	result := &specs.Property{
		Path:     path,
		Name:     descriptor.GetName(),
		Comment:  descriptor.GetSourceInfo().GetLeadingComments(),
		Position: descriptor.GetNumber(),
		Type:     Types[descriptor.GetType()],
		Label:    Labels[descriptor.GetLabel()],
		Options:  specs.Options{},
	}

	if descriptor.GetType() != protobuf.FieldDescriptorProto_TYPE_MESSAGE {
		return result
	}

	fields := descriptor.GetMessageType().GetFields()
	result.Nested = make(map[string]*specs.Property, len(fields))

	for _, field := range fields {
		result.Nested[field.GetName()] = NewProperty(template.JoinPath(path, field.GetName()), field)
	}

	return result
}
