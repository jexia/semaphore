package sprintf

import (
	"fmt"
	"strings"

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

	return FormatWithFunc(strtoa)(precision), nil
}

func strtoa(precision Precision, value interface{}) (string, error) {
	if value == nil {
		return "", errNoValue
	}

	casted, ok := value.(string)
	if !ok {
		return "", errNonStringType
	}

	if precision.Width == 0 {
		return casted, nil
	}

	var builder strings.Builder
	for i := int64(0); i < precision.Width; i++ {
		builder.WriteByte(casted[i])
	}

	return builder.String(), nil
}
