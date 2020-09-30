package sprintf

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
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
func FormatJSON(store references.Store, argument *specs.Property) (string, error) {
	var property = encoder{refs: store, property: argument}

	bb, err := json.Marshal(property)
	if err != nil {
		return "", err
	}

	return string(bb), nil
}

type encoder struct {
	property *specs.Property
	refs     references.Store
}

func (enc encoder) MarshalJSON() ([]byte, error) {
	switch {
	case enc.property.Repeated != nil:
		return repeated{property: enc.property, refs: enc.refs}.MarshalJSON()
	case enc.property.Message != nil:
		return message{property: enc.property, refs: enc.refs}.MarshalJSON()
	case enc.property.Enum != nil && enc.property.Reference != nil:
		reference := enc.refs.Load(enc.property.Reference.Resource, enc.property.Reference.Path)
		if reference == nil {
			return null, nil
		}

		var enum = enc.property.Enum.Positions[*reference.Enum]
		if enum == nil {
			return json.Marshal(*reference.Enum)
		}

		return json.Marshal(enum.Key)
	case enc.property.Scalar != nil:
		value := enc.property.Scalar.Default

		var reference = enc.refs.Load(enc.property.Reference.Resource, enc.property.Reference.Path)
		if reference != nil {
			value = reference.Value
		}

		return json.Marshal(value)
	}

	return null, nil
}

type repeated struct {
	property *specs.Property
	refs     references.Store
}

func (r repeated) MarshalJSON() ([]byte, error) {
	if r.property.Reference == nil {
		return null, nil
	}

	var reference = r.refs.Load(r.property.Reference.Resource, r.property.Reference.Path)
	if reference == nil || reference.Repeated == nil {
		return null, nil
	}

	var buff = bytes.NewBufferString("[")

	for index, store := range reference.Repeated {
		if index > 0 {
			buff.WriteString(",")
		}

		item, err := r.property.Repeated.Template()
		if err != nil {
			return nil, fmt.Errorf("failed to encode repeated item: %w", err)
		}

		// the item template does include its own reference
		item.Reference = &specs.PropertyReference{}

		bb, err := encoder{property: &specs.Property{Template: item}, refs: store}.MarshalJSON()
		if err != nil {
			return nil, err
		}

		buff.Write(bb)
	}

	buff.WriteString("]")

	return buff.Bytes(), nil
}

type message struct {
	property *specs.Property
	refs     references.Store
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
		bb, err := (&encoder{property: prop, refs: m.refs}).MarshalJSON()
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
