package sprintf

import (
	"fmt"
	"strconv"

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

	return FormatWithFunc(itoa)(precision), nil
}

func itoa(_ Precision, value interface{}) (string, error) {
	switch t := value.(type) {
	case nil:
		return "", errNoValue
	case int:
		return strconv.FormatInt(int64(t), 10), nil
	case int32:
		return strconv.FormatInt(int64(t), 10), nil
	case int64:
		return strconv.FormatInt(t, 10), nil
	case uint:
		return strconv.FormatUint(uint64(t), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(t), 10), nil
	case uint64:
		return strconv.FormatUint(t, 10), nil
	default:
		return "", errNonIntegerType
	}
}
