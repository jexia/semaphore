package protobuffers

import (
	"fmt"
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

	for index, descriptor := range descriptors {
		log.Printf("DESCRIPTOR N%d", index)

		for _, message := range descriptor.GetMessageTypes() {
			result[message.GetFullyQualifiedName()] = NewMessage("", registry, message)
		}
	}

	log.Println("SCHEMA IS PARESED")

	// log.Printf("%s", result["semaphore.greeter.Response"].Message["meta"].Message["child"].Message["child"].Message)

	return result
}

// NewMessage constructs a schema Property with the given message descriptor
func NewMessage(path string, registry map[string]*specs.Property, descriptor *desc.MessageDescriptor) *specs.Property {
	log.Println("TOP:", path, descriptor.GetFullyQualifiedName())

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

func AddProperty(registry, messages map[string]*specs.Property, path string, descriptor *desc.FieldDescriptor) bool {
	var (
		id = fmt.Sprintf("%d", descriptor.AsFieldDescriptorProto().Number)
		//descriptor.GetFullyQualifiedName()
		name = descriptor.GetName()
	)

	log.Println("FQN", id, descriptor.GetFullyQualifiedName())
	// log.Println("number", descriptor.AsFieldDescriptorProto().Number)
	log.Println()

	property, ok := registry[id]
	if ok {
		messages[name] = property

		return true
	}

	// log.Println("FQN", id, descriptor.GetFullyQualifiedName())
	// // log.Println("number", descriptor.AsFieldDescriptorProto().Number)
	// log.Println()

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
			// NewMessage(template.JoinPath(path, field.GetName()), registry, field)

			if AddProperty(registry, property.Message, template.JoinPath(path, field.GetName()), field) {
				return true
			}
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

	return false
}
