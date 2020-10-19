package json

import (
	"github.com/francoispqt/gojay"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

// Object represents a JSON object
type Object struct {
	resource string
	message  specs.Message
	store    references.Store
	length   int
}

// NewObject constructs a new object encoder/decoder for the given specs
func NewObject(resource string, message specs.Message, store references.Store) *Object {
	return &Object{
		resource: resource,
		message:  message,
		store:    store,
		length:   len(message),
	}
}

// MarshalJSONObject encodes the given specs object into the given gojay encoder
func (object *Object) MarshalJSONObject(encoder *gojay.Encoder) {
	for _, prop := range object.message.SortedProperties() {
		encodeElementKey(encoder, object.resource, prop.Name, prop.Template, object.store)
	}
}

// UnmarshalJSONObject unmarshals the given specs into the configured reference store
func (object *Object) UnmarshalJSONObject(decoder *gojay.Decoder, key string) error {
	if object == nil {
		return nil
	}

	property, has := object.message[key]
	if !has {
		return nil
	}

	return decodeElement(decoder, object.resource, property.Path, property.Template, object.store)
}

// NKeys returns the amount of available keys inside the given object
func (object *Object) NKeys() int {
	return object.length
}

// IsNil returns whether the given object is null or not
func (object *Object) IsNil() bool {
	return object == nil
}
