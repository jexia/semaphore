package proto

import (
	"github.com/jexia/maestro/schema/protoc"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/types"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

// NewMessage attempts to construct a new proto message descriptor for the given specs property
func NewMessage(resource string, specs map[string]*specs.Property) (*desc.MessageDescriptor, error) {
	msg := builder.NewMessage(resource)
	err := ConstructMessage(msg, specs)
	if err != nil {
		return nil, err
	}

	return msg.Build()
}

// ConstructMessage constructs a proto message of the given specs into the given message builders
func ConstructMessage(msg *builder.MessageBuilder, specs map[string]*specs.Property) (err error) {
	for key, prop := range specs {
		if prop.Type == types.Message {
			nested := builder.NewMessage(key)
			err = ConstructMessage(nested, prop.Nested)
			if err != nil {
				return err
			}

			message := builder.FieldTypeMessage(nested)
			field := builder.NewField(key, message)
			field.SetJsonName(key)
			field.SetLabel(protoc.ProtoLabels[prop.Label])
			field.SetComments(builder.Comments{
				LeadingComment: prop.Desciptor.GetComment(),
			})

			err = msg.TryAddField(field.SetNumber(prop.Desciptor.GetPosition()))
			if err != nil {
				return err
			}

			continue
		}

		typ := builder.FieldTypeScalar(protoc.ProtoTypes[prop.Type])
		field := builder.NewField(key, typ)
		field.SetJsonName(key)
		field.SetLabel(protoc.ProtoLabels[prop.Label])
		field.SetComments(builder.Comments{
			LeadingComment: prop.Desciptor.GetComment(),
		})

		err = msg.TryAddField(field.SetNumber(prop.Desciptor.GetPosition()))
		if err != nil {
			return err
		}
	}

	return nil
}
