package sprintf

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

var scanner Scanner

func init() {
	var (
		defaultFormatDetector = NewRadix()

		formatters = []Constructor{
			String{},
			JSON{},
			Float{},
			Int{},
		}
	)

	for _, formatter := range formatters {
		if err := defaultFormatDetector.Register(formatter); err != nil {
			panic(err)
		}
	}

	scanner = NewScanner(defaultFormatDetector)
}

// Formatter is a function to be called in order to format the argument value.
type Formatter func(store references.Store, argument *specs.Property) (string, error)

// TypeChecker determines if the provided type can be formatted by the formatter.
type TypeChecker interface {
	CanFormat(types.Type) bool
}

// Constructor is a formatter constructor.
type Constructor interface {
	fmt.Stringer

	TypeChecker

	Formatter(Precision) (Formatter, error)
}

// ValueFormatter formats a scalar value with provided precision.
type ValueFormatter func(Precision, interface{}) (string, error)

// FormatWithFunc does common operations (e.g. retrieves scalar value from the reference store).
func FormatWithFunc(valueFormatter ValueFormatter) func(Precision) Formatter {
	return func(precision Precision) Formatter {
		return func(store references.Store, argument *specs.Property) (string, error) {
			var value interface{}

			if argument.Scalar != nil {
				value = argument.Scalar.Default
			}

			if argument.Reference != nil {
				if ref := store.Load(argument.Reference.Resource, argument.Reference.Path); ref != nil {
					value = ref.Value
				}
			}

			return valueFormatter(precision, value)
		}
	}
}
