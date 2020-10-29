package json

import (
	"log"

	"github.com/francoispqt/gojay"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

// Array represents a JSON array.
type Array struct {
	path     string
	template specs.Template
	repeated specs.Template
	tracker  references.Tracker
	store    references.Store
}

// NewArray creates a new array to be JSON encoded/decoded.
func NewArray(path string, template specs.Template, store references.Store, tracker references.Tracker) *Array {
	log.Println("T REF:", template.Reference)

	// TODO: find a better implementation/name
	combi, err := template.Repeated.Template()
	if err != nil {
		panic(err)
	}

	log.Println("C REF:", combi.Reference)

	return &Array{
		path:     path,
		template: combi,
		repeated: template,
		store:    store,
		tracker:  tracker,
	}
}

// MarshalJSONArray encodes the array into the given gojay encoder.
func (array *Array) MarshalJSONArray(encoder *gojay.Encoder) {
	if array.repeated.Reference == nil {
		ptrack := array.tracker.Resolve(array.path)
		array.tracker.Track(ptrack, 0)

		for _, template := range array.repeated.Repeated {
			encode(encoder, array.path, template, array.store, array.tracker)
			array.tracker.Next(ptrack)
		}

		return
	}

	length := array.store.Length(array.tracker.Resolve(array.repeated.Reference.String()))
	if length == 0 {
		return
	}

	ptrack := array.tracker.Resolve(array.path)
	array.tracker.Track(ptrack, 0)

	for index := 0; index < length; index++ {
		encode(encoder, array.path, array.template, array.store, array.tracker)
		array.tracker.Next(ptrack)
	}
}

// UnmarshalJSONArray unmarshals the given specs into the configured reference store.
func (array *Array) UnmarshalJSONArray(decoder *gojay.Decoder) error {
	array.tracker.Track(array.path, decoder.Index())
	array.store.Define(array.tracker.Resolve(array.path), decoder.Index()+1) // assuming that decoder increments by one
	return decode(decoder, array.path, array.template, array.store, array.tracker)
}

// IsNil returns whether the given array is null or not.
func (array *Array) IsNil() bool {
	return array == nil
}
