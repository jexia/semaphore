package sprintf

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/jexia/semaphore/v2/pkg/references"
	"github.com/jexia/semaphore/v2/pkg/specs"
	"github.com/jexia/semaphore/v2/pkg/specs/types"
)

var null = []byte("null")

// JSON formatter.
type JSON struct{}

func (JSON) String() string { return "json" }

// CanFormat checks whether formatter accepts provided data type or not.
func (JSON) CanFormat(dataType types.Type) bool { return true }

// Formatter validates the presision and returns a JSON formatter.
func (json JSON) Formatter(precision Precision) (Formatter, error) {
	if precision.Width != 0 || precision.Scale != 0 {
		return nil, fmt.Errorf("%s formatter does not support precision", json)
	}

	return FormatJSON, nil
}

// FormatJSON prints provided argument in a JSON format.
func FormatJSON(store references.Store, tracker references.Tracker, argument *specs.Property) (string, error) {
	property := encoder{store: store, tracker: tracker, property: argument}

	bb, err := json.Marshal(property)
	if err != nil {
		return "", err
	}

	return string(bb), nil
}

type encoder struct {
	property *specs.Property
	store    references.Store
	tracker  references.Tracker
}

func (enc encoder) MarshalJSON() ([]byte, error) {
	switch {
	case enc.property.Repeated != nil:
		return repeated{property: enc.property, store: enc.store, tracker: enc.tracker}.MarshalJSON()
	case enc.property.Message != nil:
		return message{property: enc.property, store: enc.store, tracker: enc.tracker}.MarshalJSON()
	case enc.property.Enum != nil && enc.property.Reference != nil:
		reference := enc.store.Load(enc.tracker.Resolve(enc.property.Reference.String()))
		if reference == nil {
			return null, nil
		}

		enum := enc.property.Enum.Positions[*reference.Enum]
		if enum == nil {
			return json.Marshal(*reference.Enum)
		}

		return json.Marshal(enum.Key)
	case enc.property.Scalar != nil:
		value := enc.property.Scalar.Default

		reference := enc.store.Load(enc.tracker.Resolve(enc.property.Reference.String()))
		if reference != nil {
			value = reference.Value
		}

		return json.Marshal(value)
	}

	return null, nil
}

type repeated struct {
	property *specs.Property
	store    references.Store
	tracker  references.Tracker
}

func (r repeated) MarshalJSON() ([]byte, error) {
	if r.property.Reference == nil {
		return null, nil
	}

	path := r.tracker.Resolve(r.property.Reference.String())
	r.tracker.Track(path, 0)
	length := r.store.Length(path)
	if length == 0 {
		return null, nil
	}

	buff := bytes.NewBufferString("[")
	item, err := r.property.Repeated.Template()
	if err != nil {
		return nil, fmt.Errorf("failed to encode repeated item: %w", err)
	}

	for index := 0; index < length; index++ {
		if index > 0 {
			buff.WriteString(",")
		}

		bb, err := encoder{property: &specs.Property{Template: item}, store: r.store, tracker: r.tracker}.MarshalJSON()
		if err != nil {
			return nil, err
		}

		buff.Write(bb)
		r.tracker.Next(path)
	}

	buff.WriteString("]")

	return buff.Bytes(), nil
}

type message struct {
	property *specs.Property
	store    references.Store
	tracker  references.Tracker
}

func (m message) MarshalJSON() ([]byte, error) {
	if m.property.Message == nil {
		return null, nil
	}

	var (
		buff     = bytes.NewBufferString("{")
		firstKey = true
	)

	for key, prop := range m.property.Message {
		bb, err := (&encoder{property: prop, store: m.store, tracker: m.tracker}).MarshalJSON()
		if err != nil {
			return nil, err
		}

		if !firstKey {
			buff.WriteString(",")
		}

		buff.WriteString(`"` + key + `":`)
		buff.Write(bb)

		firstKey = false
	}

	buff.WriteString("}")

	return buff.Bytes(), nil
}
