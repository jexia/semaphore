package sprintf

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

// Int formatter.
type Int struct{}

func (Int) String() string { return "d" }

// CanFormat checks whether formatter accepts provided data type or not.
func (Int) CanFormat(dataType types.Type) bool {
	switch dataType {
	case types.Int32, types.Int64, types.Uint32, types.Uint64:
		return true
	default:
		return false
	}
}

// Formatter creates new float formatter.
func (i Int) Formatter(precision Precision) (Formatter, error) {
	if precision.Width != 0 && precision.Scale != 0 {
		return nil, fmt.Errorf("%q formatter does not support precision", i)
	}

	return FormatInt, nil
}

// FormatInt prints provided argument as an integer.
func FormatInt(store references.Store, argument *specs.Property) (string, error) {
	var value interface{}

	if argument.Default != nil {
		value = argument.Default
	}

	if argument.Reference != nil {
		if ref := store.Load(argument.Reference.Resource, argument.Reference.Path); ref != nil {
			value = ref.Value
		}
	}

	return itoa(value)
}

func itoa(value interface{}) (string, error) {
	switch t := value.(type) {
	case nil:
		return "", errNoValue
	case int:
		return strconv.FormatInt(int64(t), 10), nil
	case int32:
		return strconv.FormatInt(int64(t), 10), nil
	case int64:
		return strconv.FormatInt(t, 10), nil
	default:
		return "", errors.New("not an integer")
	}
}
