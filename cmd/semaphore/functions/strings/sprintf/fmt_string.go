package sprintf

import (
	"errors"
	"fmt"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

// String formatter.
type String struct{}

func (String) String() string { return "s" }

// CanFormat checks whether formatter accepts provided data type or not.
func (String) CanFormat(dataType types.Type) bool {
	switch dataType {
	case types.String:
		return true
	default:
		return false
	}
}

// Formatter creates new string formatter.
func (str String) Formatter(precision Precision) (Formatter, error) {
	if precision.Width != 0 && precision.Scale != 0 {
		return nil, fmt.Errorf("%q formatter does not support precision", str)
	}

	return FormatString, nil
}

// FormatString prints provided argument as a string.
func FormatString(store references.Store, argument *specs.Property) (string, error) {
	var value interface{}

	if argument.Default != nil {
		value = argument.Default
	}

	if argument.Reference != nil {
		if ref := store.Load(argument.Reference.Resource, argument.Reference.Path); ref != nil {
			value = ref.Value
		}
	}

	if value == nil {
		return "", nil
	}

	casted, ok := value.(string)
	if !ok {
		return "", errors.New("not a string")
	}

	return casted, nil
}
