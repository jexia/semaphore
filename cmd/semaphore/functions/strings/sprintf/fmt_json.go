package sprintf

import (
	"fmt"
	"strings"

	"github.com/francoispqt/gojay"
	"github.com/jexia/semaphore/pkg/codec/json"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

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
	var (
		builder  strings.Builder
		property = &encoder{
			resource: "",
			refs:     store,
			property: argument,
		}
		encoder = gojay.NewEncoder(&builder)
	)

	if err := encoder.Encode(property); err != nil {
		return "", err
	}

	return builder.String(), nil
}

type encoder struct {
	resource string
	property *specs.Property
	refs     references.Store
}

func (enc *encoder) MarshalJSONObject(encoder *gojay.Encoder) {
	json.MarshalProperty(enc.refs, enc.resource, enc.property, encoder)
}

func (enc *encoder) IsNil() bool {
	return enc.property == nil
}
