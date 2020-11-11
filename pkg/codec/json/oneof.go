package json

import (
	"errors"

	"github.com/francoispqt/gojay"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
)

// OneOf represents a JSON object
type OneOf struct {
	path     string
	template specs.Template
	store    references.Store
	tracker  references.Tracker
	isSet    bool
}

// NewOneOf constructs a new object encoder/decoder for the given specs
func NewOneOf(path string, template specs.Template, store references.Store, tracker references.Tracker) *OneOf {
	return &OneOf{
		path:     path,
		template: template,
		store:    store,
		tracker:  tracker,
	}
}

// MarshalJSONObject encodes the given specs object into the given gojay encoder
func (oneOf *OneOf) MarshalJSONObject(encoder *gojay.Encoder) {
	for _, prop := range oneOf.template.OneOf {
		encodeKey(encoder, template.JoinPath(oneOf.path, prop.Name), prop.Name, prop.Template, oneOf.store, oneOf.tracker)
	}
}

// UnmarshalJSONObject unmarshals the given specs into the configured reference store
func (oneOf *OneOf) UnmarshalJSONObject(decoder *gojay.Decoder, key string) error {
	if oneOf.template.OneOf == nil {
		return nil
	}

	if oneOf.isSet {
		return errors.New("only a single field is allowed to be set")
	}

	oneOf.isSet = true

	property, has := oneOf.template.OneOf[key]
	if !has {
		return nil
	}

	oneOf.store.Define(oneOf.path, len(oneOf.template.Message))

	return decode(decoder, template.JoinPath(oneOf.path, key), property.Template, oneOf.store, oneOf.tracker)
}

// NKeys returns the amount of available keys inside the given oneof
func (oneOf *OneOf) NKeys() int {
	if oneOf.template.OneOf == nil {
		return 0
	}

	return len(oneOf.template.OneOf)
}

// IsNil returns whether the given oneof is null or not
func (oneOf *OneOf) IsNil() bool {
	return oneOf == nil
}
