package sprintf

import (
	"fmt"

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
	panic("not implemented")
}
