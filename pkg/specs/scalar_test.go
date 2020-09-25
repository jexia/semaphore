package specs

import (
	"testing"

	"github.com/jexia/semaphore/pkg/specs/types"
)

func TestScalarUnmarshalInvalidJSON(t *testing.T) {
	payload := "non json string"
	prop := Property{
		Template: Template{
			Scalar: &Scalar{},
		},
	}

	err := prop.Scalar.UnmarshalJSON([]byte(payload))
	if err == nil {
		t.Error("expected error got nil")
	}
}

func TestScalarUnmarshalNil(t *testing.T) {
	var scalar *Scalar
	err := scalar.UnmarshalJSON(nil)
	if err != nil {
		t.Error(err)
	}
}

func TestScalarCompare(t *testing.T) {
	shouldBeEqual := func(t *testing.T, one, another *Scalar) {
		if err := one.Compare(another); err != nil {
			t.Fatalf("CompareScalars() returns unexpected error: %v", err)
		}
	}

	shouldNotBeEqual := func(t *testing.T, one, another *Scalar) {
		if one.Compare(another) == nil {
			t.Fatalf("CompareScalars() = nil, want an error")
		}
	}

	t.Run("should be equal", func(t *testing.T) {
		shouldBeEqual(t, &Scalar{Type: types.Int32}, &Scalar{Type: types.Int32})
	})

	t.Run("nil scalars should be equal", func(t *testing.T) {
		shouldBeEqual(t, nil, nil)
	})

	t.Run("should not be equal", func(t *testing.T) {
		shouldNotBeEqual(t, &Scalar{Type: types.Int32}, &Scalar{Type: types.String})
	})

	t.Run("nil and not nil scalars should not be equal", func(t *testing.T) {
		shouldNotBeEqual(t, &Scalar{}, nil)
		shouldNotBeEqual(t, nil, &Scalar{})
	})
}
