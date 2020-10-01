package xml

import (
	"encoding/xml"
	"errors"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

// Array represents an array of values/references.
type Array struct {
	resource  string
	prefix    string
	name      string
	template  specs.Template
	repeated  specs.Repeated
	reference *specs.PropertyReference
	store     references.Store
}

// NewArray constructs a new XML array encoder/decoder.
func NewArray(resource, prefix, name string, template specs.Template, repeated specs.Repeated, reference *specs.PropertyReference, store references.Store) *Array {
	return &Array{
		resource:  resource,
		prefix:    prefix,
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
		return nil
	}

	var reference = array.store.Load(array.reference.Resource, array.reference.Path)
	if reference == nil {
		return nil
	}

	if reference.Repeated == nil {
		return errors.New("reference does not contain repeated value")
	}

	for _, store := range reference.Repeated {
		var template = array.template.Clone()
		template.Reference = new(specs.PropertyReference)

		if err := encodeElement(encoder, array.name, template, store); err != nil {
			return err
		}
	}

	return nil
}

// UnmarshalXML decodes XML input into the receiver of type specs.Repeated.
func (array *Array) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var (
		store     = references.NewReferenceStore(1)
		reference = array.store.Load(array.resource, buildPath(array.prefix, array.name))
	)

	if reference == nil {
		reference = &references.Reference{
			Path: buildPath(array.prefix, array.name),
		}

		array.store.StoreReference(array.resource, reference)
	}

	// TODO: fixme

	if err := decodeElement(decoder, start, "", "", "", array.template, store); err != nil {
		return err
	}

	reference.Repeated = append(reference.Repeated, store)

	return nil
}
