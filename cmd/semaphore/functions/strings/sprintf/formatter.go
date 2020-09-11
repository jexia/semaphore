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

type TypeChecker interface {
	CanFormat(types.Type) bool
}

// Constructor is a formatter constructor.
type Constructor interface {
	fmt.Stringer

	TypeChecker

	Formatter(Precision) (Formatter, error)
}
