package xml

import (
	"encoding/xml"

	"github.com/jexia/semaphore/v2/pkg/references"
	"github.com/jexia/semaphore/v2/pkg/specs"
	"github.com/jexia/semaphore/v2/pkg/specs/template"
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
	// TODO: find a better implementation/name
	combi, err := template.Repeated.Template()
	if err != nil {
		panic(err)
	}

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
			if err := encodeElement(
				encoder,
				array.name,
				array.path,
				item,
				array.store,
				array.tracker,
			); err != nil {
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
		if err := encodeElement(
			encoder,
			array.name,
			template.JoinPath(array.path, array.name),
			array.template,
			array.store,
			array.tracker,
		); err != nil {
			return err
		}

		array.tracker.Next(ptrack)
	}

	return nil
}

// UnmarshalXML decodes XML input into the receiver of type specs.Repeated.
func (array *Array) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var (
		path  = template.JoinPath(array.path, array.name)
		index = array.store.Length(array.tracker.Resolve(path))
	)

	array.tracker.Track(path, index)
	array.store.Define(array.tracker.Resolve(path), index+1)

	return decodeElement(
		decoder,
		start,
		array.name,
		array.path,
		array.template,
		array.store,
		array.tracker,
	)
}
