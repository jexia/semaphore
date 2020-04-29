package proto

import (
	"github.com/jexia/maestro/pkg/definitions/protoc"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/types"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

// NewMessage attempts to construct a new proto message descriptor for the given specs property
func NewMessage(resource string, specs map[string]*specs.Property) (*desc.MessageDescriptor, error) {
	msg := builder.NewMessage(resource)
	err := ConstructMessage(msg, nil, specs)
	if err != nil {
		return nil, err
	}

	return msg.Build()
}

// ConstructMessage constructs a proto message of the given specs into the given message builders
func ConstructMessage(msg *builder.MessageBuilder, file *builder.FileBuilder, specs map[string]*specs.Property) (err error) {
	for key, prop := range specs {
		if prop.Type == types.Message {
			name := key // TODO: the name is not unique causing that two properties with the same name will conflict
			nested := builder.NewMessage(name)
			err = ConstructMessage(nested, file, prop.Nested)
			if err != nil {
				return err
			}

			message := builder.FieldTypeMessage(nested)
			field := builder.NewField(name, message)
			field.SetJsonName(name)
			field.SetLabel(protoc.ProtoLabels[prop.Label])
			field.SetComments(builder.Comments{
				LeadingComment: prop.Comment,
			})

			err = msg.TryAddField(field.SetNumber(prop.Position))
			if err != nil {
				return err
			}

			if file != nil {
				err = file.TryAddMessage(nested)
				if err != nil {
					return err
				}
			}

			continue
		}

		typ := builder.FieldTypeScalar(protoc.ProtoTypes[prop.Type])
		field := builder.NewField(key, typ)
		field.SetJsonName(key)
		field.SetLabel(protoc.ProtoLabels[prop.Label])
		field.SetComments(builder.Comments{
			LeadingComment: prop.Comment,
		})

		err = msg.TryAddField(field.SetNumber(prop.Position))
		if err != nil {
			return err
		}
	}

	return nil
}
