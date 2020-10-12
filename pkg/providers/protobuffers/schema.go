package protobuffers

import (
	"log"

	protobuf "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/template"
	"github.com/jhump/protoreflect/desc"
)

// NewSchema constructs a new schema manifest from the given file descriptors
func NewSchema(descriptors []*desc.FileDescriptor) specs.Schemas {
	log.Println("PARESING THE SCHEMA")

	result := make(specs.Schemas)

	registry := make(map[string]*specs.Property)

	for _, descriptor := range descriptors {
		for _, message := range descriptor.GetMessageTypes() {
			result[message.GetFullyQualifiedName()] = NewMessage("", registry, message)
		}
	}

	log.Println("SCHEMA IS PARESED")

	return result
}

// NewMessage constructs a schema Property with the given message descriptor
func NewMessage(path string, registry map[string]*specs.Property, descriptor *desc.MessageDescriptor) *specs.Property {
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
		AddProperty(registry, result.Message, template.JoinPath(path, field.GetName()), field)
	}

	return result
}

func AddProperty(registry, messages map[string]*specs.Property, path string, descriptor *desc.FieldDescriptor) {
	var (
		id   = descriptor.GetFullyQualifiedName()
		name = descriptor.GetName()
	)

	log.Println(id, name)

	property, ok := registry[id]
	if ok {
		messages[name] = property

		return
	}

	property = &specs.Property{
		Path:        path,
		Name:        name,
		Description: descriptor.GetSourceInfo().GetLeadingComments(),
		Position:    descriptor.GetNumber(),
		Options:     specs.Options{},
		Label:       Labels[descriptor.GetLabel()],
		Template: specs.Template{
			Identifier: id,
		},
	}

	registry[id] = property
	messages[name] = property

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

		property.Enum = &specs.Enum{
			Name:        enum.GetName(),
			Description: enum.GetSourceInfo().GetLeadingComments(),
			Keys:        keys,
			Positions:   positions,
		}
	case descriptor.GetType() == protobuf.FieldDescriptorProto_TYPE_MESSAGE:
		var fields = descriptor.GetMessageType().GetFields()
		property.Message = make(specs.Message, len(fields))

		for _, field := range fields {
			AddProperty(registry, property.Message, template.JoinPath(path, field.GetName()), field)
		}
	default:
		property.Scalar = &specs.Scalar{
			Type: Types[descriptor.GetType()],
		}
	}

	if descriptor.GetLabel() == protobuf.FieldDescriptorProto_LABEL_REPEATED {
		property.Label = labels.Optional
		property.Template = specs.Template{
			Repeated: specs.Repeated{property.Template},
		}
	}
}

// NewProperty constructs a schema Property with the given field descriptor
func NewProperty(registry map[string]*specs.Property, path string, descriptor *desc.FieldDescriptor) *specs.Property {
	log.Println(descriptor.GetFullyQualifiedName())

	property, ok := registry[descriptor.GetFullyQualifiedName()]
	if ok {
		return property
	}

	property = &specs.Property{
		Path:        path,
		Name:        descriptor.GetName(),
		Description: descriptor.GetSourceInfo().GetLeadingComments(),
		Position:    descriptor.GetNumber(),
		Options:     specs.Options{},
		Label:       Labels[descriptor.GetLabel()],
		Template: specs.Template{
			Identifier: descriptor.GetFullyQualifiedName(),
		},
	}

	switch {
	// case descriptor.GetOneOf() != nil:
	// 	var choices = descriptor.GetOneOf().GetChoices()
	// 	property.OneOf = &specs.OneOf{
	// 		Choices: make(map[string]*specs.Property, len(choices)),
	// 	}
	//
	// 	for _, choice := range choices {
	// 		// result.OneOf.Choices[choice.GetName()] = NewProperty(template.JoinPath(path, choice.GetName()), choice)
	// 		log.Println("choice:", choice)
	// 	}
	//
	// 	// // descriptor.GetOneOf().GetOneOfOptions()
	// 	// log.Printf("ONEOF: %#v\n", descriptor.GetOneOf().GetName())
	// 	// log.Printf("PARENT: %#v\n\n\n", descriptor.GetOneOf().GetParent())

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

		property.Enum = &specs.Enum{
			Name:        enum.GetName(),
			Description: enum.GetSourceInfo().GetLeadingComments(),
			Keys:        keys,
			Positions:   positions,
		}
	case descriptor.GetType() == protobuf.FieldDescriptorProto_TYPE_MESSAGE:
		var fields = descriptor.GetMessageType().GetFields()
		property.Message = make(specs.Message, len(fields))

		for _, field := range fields {
			property.Message[field.GetName()] = NewProperty(registry, template.JoinPath(path, field.GetName()), field)
		}
	default:
		property.Scalar = &specs.Scalar{
			Type: Types[descriptor.GetType()],
		}
	}

	if descriptor.GetLabel() == protobuf.FieldDescriptorProto_LABEL_REPEATED {
		property.Label = labels.Optional
		property.Template = specs.Template{
			Repeated: specs.Repeated{property.Template},
		}
	}

	return property
}
