package sprintf

import (
	"fmt"

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

		fmt.Println()
		fmt.Println(precision.String())
		fmt.Printf("%T", value)
		fmt.Println()
		fmt.Println()

		return fmt.Sprintf(precision.String()+"f", value), nil
	}
}
