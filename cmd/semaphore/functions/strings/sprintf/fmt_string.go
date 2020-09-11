package sprintf

import (
	"fmt"
	"strings"

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
	if precision.Scale != 0 {
		return nil, fmt.Errorf("%q formatter does not support scale", str)
	}

	return FormatString(precision.Width), nil
}

// FormatString prints provided argument as a string.
func FormatString(length int64) Formatter {
	return func(store references.Store, argument *specs.Property) (string, error) {
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
			return "", errNonStringType
		}

		if length == 0 {
			return casted, nil
		}

		var builder strings.Builder
		for i := int64(0); i < length; i++ {
			if err := builder.WriteByte(casted[i]); err != nil {
				return "", err
			}
		}

		return builder.String(), nil
	}
}
