package avros

import (
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/template"
	"github.com/jexia/semaphore/pkg/specs/types"
)

// AvroSchema impliments avro high level schema
type AvroSchema struct {
	Type      string        `json:"type"`
	Name      string        `json:"name"`
	Namespace string        `json:"namespace"`
	Symbols   []string      `json:"symbols"`
	Fields    []*AvroSchema `json:"fields"`
}

// NewSchema constructs a new schema manifest from the given avro schema
func NewSchema(descriptors []*AvroSchema) specs.Schemas {
	result := make(specs.Schemas, 0)

	for _, desc := range descriptors {
		result[desc.Namespace] = NewMessage("", desc)
	}

	return result
}

// NewMessage constructs a schema Property with the given avro schema
func NewMessage(path string, message *AvroSchema) *specs.Property {
	result := &specs.Property{
		Path:     message.Namespace,
		Name:     message.Name,
		Position: 1,
		Label:    labels.Optional,
		Template: &specs.Template{
			Message: make(specs.Message, len(message.Fields)),
		},
		Options: specs.Options{},
	}

	for _, field := range message.Fields {
		result.Message[field.Name] = NewProperty(template.JoinPath(message.Namespace, message.Name, field.Name), field)
	}

	return result
}

// NewProperty constructs a schema Property with the given avro schema
func NewProperty(path string, message *AvroSchema) *specs.Property {
	result := &specs.Property{
		Path:    path,
		Name:    message.Name,
		Options: specs.Options{},
	}

	switch message.Type {
	case AvroTypes[types.Enum]:
		keys := map[string]*specs.EnumValue{}
		positions := map[int32]*specs.EnumValue{}
		for i, value := range message.Symbols {
			result := &specs.EnumValue{
				Key:      value,
				Position: int32(i),
			}
			keys[value] = result
			positions[int32(i)] = result
		}

		result.Enum = &specs.Enum{
			Name:      message.Name,
			Keys:      keys,
			Positions: positions,
		}
	case AvroTypes[types.Message]:
		fields := message.Fields
		result.Message = make(specs.Message, len(fields))
		for _, field := range fields {
			result.Message[field.Name] = NewProperty(template.JoinPath(path, message.Name, field.Name), field)
		}
		break
	default:
		result.Scalar = &specs.Scalar{
			Type: Types[message.Type],
		}
	}

	return result
}
