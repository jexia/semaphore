package json

import (
	"github.com/francoispqt/gojay"
	"github.com/jexia/semaphore/v2/pkg/references"
	"github.com/jexia/semaphore/v2/pkg/specs"
)

// Object represents a JSON object
type Object struct {
	path     string
	template specs.Template
	store    references.Store
	tracker  references.Tracker
}

// NewObject constructs a new object encoder/decoder for the given specs
func NewObject(path string, template specs.Template, store references.Store, tracker references.Tracker) *Object {
	return &Object{
		path:     path,
		template: template,
		store:    store,
		tracker:  tracker,
	}
}

// MarshalJSONObject encodes the given specs object into the given gojay encoder
func (object *Object) MarshalJSONObject(encoder *gojay.Encoder) {
	for _, prop := range object.template.Message.SortedProperties() {
		encodeKey(encoder, specs.JoinPath(object.path, prop.Name), prop.Name, prop.Template, object.store, object.tracker)
	}
}

// UnmarshalJSONObject unmarshals the given specs into the configured reference store
func (object *Object) UnmarshalJSONObject(decoder *gojay.Decoder, key string) error {
	if object.template.Message == nil {
		return nil
	}

	property, has := object.template.Message[key]
	if !has {
		return nil
	}

	object.store.Define(object.path, len(object.template.Message))

	return decode(decoder, specs.JoinPath(object.path, key), property.Template, object.store, object.tracker)
}

// NKeys returns the amount of available keys inside the given object
func (object *Object) NKeys() int {
	if object.template.Message == nil {
		return 0
	}

	return len(object.template.Message)
}

// IsNil returns whether the given object is null or not
func (object *Object) IsNil() bool {
	return object == nil
}
