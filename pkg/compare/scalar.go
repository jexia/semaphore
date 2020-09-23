package compare

import (
	"errors"
	"fmt"

	"github.com/jexia/semaphore/pkg/specs"
)

// CompareScalars compares the given scalar against of the expected and returns the first met difference as error.
func CompareScalars(expected, given *specs.Scalar) error {
	if expected == nil && given == nil {
		return nil
	}

	if expected == nil && given != nil {
		return errors.New("expected to be nil")
	}

	if expected != nil && given == nil {
		return fmt.Errorf("expected to be %v, got %v", expected.Type, nil)
	}

	if expected.Type != given.Type {
		return fmt.Errorf("expected to be %v, got %v", expected.Type, given.Type)
	}

	return nil
}
