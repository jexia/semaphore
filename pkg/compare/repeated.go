package compare

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/specs"
)

// CompareRepeated compares the given repeated structure against of the expected one and returns the first met difference as error.
func CompareRepeated(expected *specs.Repeated, given *specs.Repeated) error {
	err := CompareTemplate(expected.Template, given.Template)
	if err != nil {
		return fmt.Errorf("repeated item mismatch: %w", err)
	}

	return nil
}
