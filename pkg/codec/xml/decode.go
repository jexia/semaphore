package xml

import (
	"encoding/xml"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

func decodeNested(decoder *xml.Decoder, start xml.StartElement, prop *specs.Property, resource string, store references.Store, refs map[string]*references.Reference) error {
	if prop.Type != types.Message {
		return errNotAnObject
	}

	var nested = NewObject(resource, prop.Nested, store)

	return nested.startElement(decoder, start, refs)
}

func decodeRepeatedNested(decoder *xml.Decoder, start xml.StartElement, prop *specs.Property, resource string, refs map[string]*references.Reference) error {
	if prop.Type != types.Message {
		return errNotAnObject
	}

	var store = references.NewReferenceStore(1)

	ref, ok := refs[prop.Path]
	if !ok {
		ref = &references.Reference{
			Path: prop.Path,
		}

		refs[prop.Path] = ref
	}

	var nested = NewObject(resource, prop.Nested, store)
	if err := nested.startElement(decoder, start, refs); err != nil {
		return err
	}

	ref.Append(store)

	return nil
}

func decodeRepeatedValue(raw xml.CharData, prop *specs.Property, refs map[string]*references.Reference) error {
	var store = references.NewReferenceStore(1)

	ref, ok := refs[prop.Path]
	if !ok {
		ref = &references.Reference{
			Path: prop.Path,
		}

		refs[prop.Path] = ref
	}

	if prop.Type == types.Enum {
		enum, ok := prop.Enum.Keys[string(raw)]
		if !ok {
			return errUnknownEnum(raw)
		}

		store.StoreEnum("", "", enum.Position)
		ref.Append(store)

		return nil
	}

	value, err := types.DecodeFromString(string(raw), prop.Type)
	if err != nil {
		return err
	}

	store.StoreValue("", "", value)
	ref.Append(store)

	return nil
}

func decodeValue(raw xml.CharData, prop *specs.Property, resource string, store references.Store) error {
	var ref = &references.Reference{
		Path: prop.Path,
	}

	if prop.Type == types.Enum {
		enum, ok := prop.Enum.Keys[string(raw)]
		if !ok {
			return errUnknownEnum(raw)
		}

		ref.Enum = &enum.Position
		store.StoreReference(resource, ref)

		return nil
	}

	value, err := types.DecodeFromString(string(raw), prop.Type)
	if err != nil {
		return err
	}

	ref.Value = value
	store.StoreReference(resource, ref)

	return nil
}
