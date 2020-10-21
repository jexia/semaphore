package protobuffers

import (
	protobuf "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	tpl "github.com/jexia/semaphore/pkg/specs/template"
	"github.com/jhump/protoreflect/desc"
)

// NewSchema constructs a new schema manifest from the given file descriptors
func NewSchema(descriptors []*desc.FileDescriptor) specs.Schemas {
	var result = make(specs.Schemas)

	for _, descriptor := range descriptors {
		for _, message := range descriptor.GetMessageTypes() {
			result[message.GetFullyQualifiedName()] = NewMessage("", message)
		}
	}

	return result
}

// NewMessage constructs a schema Property with the given message descriptor
func NewMessage(path string, descriptor *desc.MessageDescriptor) *specs.Property {
	return newMessage(make(map[string]*specs.Property), make(map[string]*specs.Template), path, descriptor)
}

func newMessage(seenProperties map[string]*specs.Property, seenTemplates map[string]*specs.Template, path string, descriptor *desc.MessageDescriptor) *specs.Property {
	var (
		fields = descriptor.GetFields()
		result = &specs.Property{
			Path:        path,
			Name:        descriptor.GetFullyQualifiedName(),
			Description: descriptor.GetSourceInfo().GetLeadingComments(),
			Position:    1,
			Label:       labels.Optional,
			Template: &specs.Template{
				Message: make(specs.Message, len(fields)),
			},
			Options: specs.Options{},
		}
	)

	for _, field := range fields {
		result.Message[field.GetName()] = GetProperty(seenProperties, seenTemplates, tpl.JoinPath(path, field.GetName()), field)
	}

	return result
}

// GetTemplate creates a new template (data type) if not available yet or returns an existing one.
func GetTemplate(seenProperties map[string]*specs.Property, seenTemplates map[string]*specs.Template, path string, descriptor *desc.FieldDescriptor) *specs.Template {
	var id = descriptor.GetFullyQualifiedName()

	if messageType := descriptor.GetMessageType(); messageType != nil {
		id = messageType.GetName()
	}

	template, ok := seenTemplates[id]
	if ok {
		return template
	}

	template = &specs.Template{
		Identifier: descriptor.GetFullyQualifiedName(),
	}

	seenTemplates[id] = template

	switch {
	case descriptor.GetType() == protobuf.FieldDescriptorProto_TYPE_ENUM:
		// TODO: implement type registry for enums
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

		template.Enum = &specs.Enum{
			Name:        enum.GetName(),
			Description: enum.GetSourceInfo().GetLeadingComments(),
			Keys:        keys,
			Positions:   positions,
		}
	case descriptor.GetType() == protobuf.FieldDescriptorProto_TYPE_MESSAGE:
		var fields = descriptor.GetMessageType().GetFields()
		template.Message = make(specs.Message, len(fields))

		for _, field := range fields {
			template.Message[field.GetName()] = GetProperty(seenProperties, seenTemplates, tpl.JoinPath(path, field.GetName()), field)
		}
	default:
		template.Scalar = &specs.Scalar{
			Type: Types[descriptor.GetType()],
		}
	}

	return template
}

// GetProperty creates a new property (if not exists) otherwise returns an existing one.
func GetProperty(seenProperties map[string]*specs.Property, seenTemplates map[string]*specs.Template, path string, descriptor *desc.FieldDescriptor) *specs.Property {
	var id = descriptor.GetFullyQualifiedName()

	property, ok := seenProperties[id]
	if ok {
		return property
	}

	property = &specs.Property{
		Path:        path,
		Name:        descriptor.GetName(),
		Identifier:  id,
		Description: descriptor.GetSourceInfo().GetLeadingComments(),
		Position:    descriptor.GetNumber(),
		Options:     specs.Options{},
		Label:       Labels[descriptor.GetLabel()],
	}

	seenProperties[id] = property

	template := GetTemplate(seenProperties, seenTemplates, path, descriptor)

	if descriptor.GetLabel() == protobuf.FieldDescriptorProto_LABEL_REPEATED {
		property.Label = labels.Optional
		property.Template = &specs.Template{
			Repeated: specs.Repeated{template},
		}

		return property
	}

	property.Template = template

	return property
}
