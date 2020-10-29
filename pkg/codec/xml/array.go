package xml

import (
	"encoding/xml"
	"log"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

// Array represents an array of values/references.
type Array struct {
	name, path string
	template   specs.Template
	repeated   specs.Template
	store      references.Store
	tracker    references.Tracker
}

// NewArray constructs a new XML array encoder/decoder.
func NewArray(name, path string, template specs.Template, store references.Store, tracker references.Tracker) *Array {
	log.Println("T REF:", template.Reference)

	// TODO: find a better implementation/name
	combi, err := template.Repeated.Template()
	if err != nil {
		panic(err)
	}

	log.Println("C REF:", combi.Reference)

	return &Array{
		name:     name,
		path:     path,
		template: combi,
		repeated: template,
		store:    store,
		tracker:  tracker,
	}
}

// MarshalXML encodes the given specs object into the provided XML encoder.
func (array *Array) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	if array.repeated.Reference == nil {
		ptrack := array.tracker.Resolve(array.path)
		array.tracker.Track(ptrack, 0)

		for _, item := range array.repeated.Repeated {
			if err := encodeElement(encoder, array.name, array.path, item, array.store, array.tracker); err != nil {
				return err
			}

			array.tracker.Next(ptrack)
		}

		return nil
	}

	length := array.store.Length(array.tracker.Resolve(array.repeated.Reference.String()))
	if length == 0 {
		return nil
	}

	ptrack := array.tracker.Resolve(array.path)
	array.tracker.Track(ptrack, 0)

	for index := 0; index < length; index++ {
		if err := encodeElement(encoder, array.name, array.path, array.template, array.store, array.tracker); err != nil {
			return err
		}

		array.tracker.Next(ptrack)
	}

	return nil
}

// UnmarshalXML decodes XML input into the receiver of type specs.Repeated.
func (array *Array) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	// var (
	// 	path      = buildPath(array.reference.Path, array.name)
	// 	store     = references.NewStore(1)
	// 	reference = array.store.Load(array.reference.Resource, path)
	// )

	// if reference == nil {
	// 	reference = &references.Reference{
	// 		Path: path,
	// 	}

	// 	array.store.StoreReference(array.reference.Resource, reference)
	// }

	// if err := decodeElement(
	// 	decoder,
	// 	start,
	// 	"", // resource
	// 	"", // path
	// 	"", // name
	// 	array.template,
	// 	store,
	// ); err != nil {
	// 	return err
	// }

	// // update the reference
	// reference.Repeated = append(reference.Repeated, store)

	return nil
}
