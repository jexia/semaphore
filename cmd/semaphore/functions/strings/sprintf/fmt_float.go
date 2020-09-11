package sprintf

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

// Float formatter.
type Float struct{}

func (Float) String() string { return "f" }

// CanFormat checks whether formatter accepts provided data type or not.
func (Float) CanFormat(dataType types.Type) bool {
	switch dataType {
	case types.Float:
		return true
	default:
		return false
	}
}

// Formatter creates new float formatter.
func (fl Float) Formatter(precision Precision) (Formatter, error) {
	return FormatFloat(precision), nil
}

// FormatFloat prints provided argument as a float.
func FormatFloat(precision Precision) Formatter {
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

		var format string

		if precision.Width > 0 {
			format = strconv.FormatInt(precision.Width, 10)
		}

		if precision.Scale > 0 {
			format += "." + strconv.FormatInt(precision.Scale, 10)
		}

		return ftoa("%"+format+"f", value)
	}
}

func ftoa(format string, value interface{}) (string, error) {
	switch t := value.(type) {
	case nil:
		return "", errNoValue
	case float32:
		return fmt.Sprintf(format, t), nil
	case float64:
		return fmt.Sprintf(format, t), nil
	default:
		return "", errors.New("not a float")
	}
}
