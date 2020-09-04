package xml

import (
	"encoding/xml"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

func decodeRepeatedValue(prop *specs.Property, raw xml.CharData, refs map[string]*references.Reference) error {
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

func decodeValue(prop *specs.Property, resource string, raw xml.CharData, store references.Store) error {
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
