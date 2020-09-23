package compare

import (
	"testing"

	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

func TestCompareScalars(t *testing.T) {
	shouldBeEqual := func(t *testing.T, expected, another *specs.Scalar) {
		got := CompareScalars(expected, another)
		if got != nil {
			t.Fatalf("CompareScalars() returns unexpected error: %v", got)
		}
	}

	shouldNotBeEqual := func(t *testing.T, expected, another *specs.Scalar) {
		if CompareScalars(expected, another) == nil {
			t.Fatalf("CompareScalars() = nil, want an error")
		}
	}

	t.Run("should be equal", func(t *testing.T) {
		shouldBeEqual(t, &specs.Scalar{Type: types.Int32}, &specs.Scalar{Type: types.Int32})
	})

	t.Run("nil scalars should be equal", func(t *testing.T) {
		shouldBeEqual(t, nil, nil)
	})

	t.Run("should not be equal", func(t *testing.T) {
		shouldNotBeEqual(t, &specs.Scalar{Type: types.Int32}, &specs.Scalar{Type: types.String})
	})

	t.Run("nil and not nil scalars should not be equal", func(t *testing.T) {
		shouldNotBeEqual(t, &specs.Scalar{}, nil)
		shouldNotBeEqual(t, nil, &specs.Scalar{})
	})
}
