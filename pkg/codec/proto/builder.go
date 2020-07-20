package proto

import (
	"github.com/jexia/semaphore/pkg/providers/protobuffers"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
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
			// TODO: appending a fixed prefix is probably not a good idea.
			nested := builder.NewMessage(key + "Nested")
			err = ConstructMessage(nested, prop.Nested)
			if err != nil {
				return err
			}

			message := builder.FieldTypeMessage(nested)
			field := builder.NewField(key, message)
			field.SetJsonName(key)
			field.SetLabel(protobuffers.ProtoLabels[prop.Label])
			field.SetComments(builder.Comments{
				LeadingComment: prop.Comment,
			})

			err = msg.TryAddField(field.SetNumber(prop.Position))
			if err != nil {
				return err
			}

			err = msg.TryAddNestedMessage(nested)
			if err != nil {
				return err
			}

			continue
		}

		var typ *builder.FieldType

		if prop.Type == types.Enum {
			enum := builder.NewEnum(key + "Enum")

			for _, value := range prop.Enum.Keys {
				eval := builder.NewEnumValue(value.Key)

				eval.SetNumber(value.Position)
				eval.SetComments(builder.Comments{
					LeadingComment: value.Description,
				})

				enum.AddValue(eval)
			}

			err = msg.TryAddNestedEnum(enum)
			if err != nil {
				return err
			}

			typ = builder.FieldTypeEnum(enum)
		}

		if typ == nil {
			typ = builder.FieldTypeScalar(protobuffers.ProtoTypes[prop.Type])
		}

		field := builder.NewField(key, typ)
		field.SetJsonName(key)
		field.SetLabel(protobuffers.ProtoLabels[prop.Label])
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
