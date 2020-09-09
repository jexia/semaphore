package sprintf

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

type Formatter interface {
	fmt.Stringer

	Format(store references.Store, argument *specs.Property) (string, error)
}
