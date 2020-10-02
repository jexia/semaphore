package xml

import (
	"encoding/xml"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

// Array represents an array of values/references.
type Array struct {
	name      string
	template  specs.Template
	repeated  specs.Repeated
	reference *specs.PropertyReference
	store     references.Store
}

// NewArray constructs a new XML array encoder/decoder.
func NewArray(name string, template specs.Template, repeated specs.Repeated, reference *specs.PropertyReference, store references.Store) *Array {
	return &Array{
		name:      name,
		template:  template,
		repeated:  repeated,
		reference: reference,
		store:     store,
	}
}

// MarshalXML encodes the given specs object into the provided XML encoder.
func (array *Array) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	if array.reference == nil {
		for _, item := range array.repeated {
			if err := encodeElement(encoder, array.name, item, array.store); err != nil {
				return err
			}
		}

		return nil
	}

	var reference = array.store.Load(array.reference.Resource, array.reference.Path)
	if reference == nil || reference.Repeated == nil {
		// ignore
		return nil
	}

	array.template.Reference = new(specs.PropertyReference)

	for _, store := range reference.Repeated {
		if err := encodeElement(encoder, array.name, array.template, store); err != nil {
			return err
		}
	}

	return nil
}

// UnmarshalXML decodes XML input into the receiver of type specs.Repeated.
func (array *Array) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var (
		path      = buildPath(array.reference.Path, array.name)
		store     = references.NewReferenceStore(1)
		reference = array.store.Load(array.reference.Resource, path)
	)

	if reference == nil {
		reference = &references.Reference{
			Path: path,
		}

		array.store.StoreReference(array.reference.Resource, reference)
	}

	if err := decodeElement(
		decoder,
		start,
		"", // resource
		"", // path
		"", // name
		array.template,
		store,
	); err != nil {
		return err
	}

	// update the reference
	reference.Repeated = append(reference.Repeated, store)

	return nil
}
