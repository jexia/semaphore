package proto

import (
	"sort"

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
	// FIXME: spec properties should have a constant index
	keys := make([]string, 0, len(specs))

	for key := range specs {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		prop := specs[key]
		if prop.Type == types.TypeMessage {
			nested := builder.NewMessage(key)
			err = ConstructMessage(nested, prop.Nested)
			if err != nil {
				return err
			}

			field := builder.FieldTypeMessage(nested)
			err = msg.TryAddField(builder.NewField(key, field))
			if err != nil {
				return err
			}

			continue
		}

		typ := builder.FieldTypeScalar(protoc.ProtoTypes[prop.Type])
		field := builder.NewField(key, typ)
		field.SetLabel(protoc.ProtoLabels[prop.Label])

		err = msg.TryAddField(field)
		if err != nil {
			return err
		}
	}

	return nil
}
