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
	result := make(specs.Schemas, 0)

	for _, descriptor := range descriptors {
		for _, message := range descriptor.GetMessageTypes() {
			result[message.GetFullyQualifiedName()] = NewMessage("", message)
		}
	}

	log.Printf("%s", result)

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
		if oneof := field.GetOneOf(); oneof != nil {
			AddOneOf(result.Message, path, field)

			continue
		}

		result.Message[field.GetName()] = NewProperty(path, field)
	}

	return result
}

// NewProperty constructs a schema Property with the given field descriptor
func NewProperty(path string, descriptor *desc.FieldDescriptor) *specs.Property {
	path = template.JoinPath(path, descriptor.GetName())

	var property = &specs.Property{
		Path:        path,
		Name:        descriptor.GetName(),
		Description: descriptor.GetSourceInfo().GetLeadingComments(),
		Position:    descriptor.GetNumber(),
		Options:     specs.Options{},
		Label:       Labels[descriptor.GetLabel()],
	}

	setTemplate(&property.Template, path, descriptor)

	return property
}

// AddOneOf constructs property of type "oneof"
func AddOneOf(message specs.Message, path string, descriptor *desc.FieldDescriptor) {
	var (
		oneOf = descriptor.GetOneOf()
		name  = oneOf.GetName()
	)

	path = template.JoinPath(path, name)

	property, ok := message[name]
	if !ok {
		property = &specs.Property{
			Path: path,
			Name: name,
			Description: func() string {
				if comments := oneOf.GetSourceInfo().LeadingComments; comments != nil {
					return *comments
				}
				return ""
			}(),
			Options: specs.Options{},
			Label:   Labels[descriptor.GetLabel()],
			Template: specs.Template{
				OneOf: make(specs.OneOf, len(oneOf.GetChoices())),
			},
		}

		message[name] = property
	}

	property.OneOf[descriptor.GetName()] = NewProperty(path, descriptor)
}

func setTemplate(template *specs.Template, path string, descriptor *desc.FieldDescriptor) {
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

		template.Enum = &specs.Enum{
			Name:        enum.GetName(),
			Description: enum.GetSourceInfo().GetLeadingComments(),
			Keys:        keys,
			Positions:   positions,
		}

		break
	case descriptor.GetType() == protobuf.FieldDescriptorProto_TYPE_MESSAGE:
		var fields = descriptor.GetMessageType().GetFields()
		template.Message = make(specs.Message, len(fields))

		for _, field := range fields {
			if oneof := field.GetOneOf(); oneof != nil {
				AddOneOf(template.Message, path, field)

				continue
			}

			template.Message[field.GetName()] = NewProperty(path, field)
		}

		break
	default:
		template.Scalar = &specs.Scalar{
			Type: Types[descriptor.GetType()],
		}
	}

	if descriptor.GetLabel() == protobuf.FieldDescriptorProto_LABEL_REPEATED {
		*template = specs.Template{
			Repeated: specs.Repeated{*template},
		}
	}
}
