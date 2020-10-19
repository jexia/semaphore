package specs

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/jexia/semaphore/pkg/specs/types"
)

func TestScalarClean(t *testing.T) {
	type test struct {
		dataType types.Type
		value    interface{}
		expected interface{}
	}

	var tests = []test{
		{dataType: types.Int32, value: uint32(42), expected: int32(42)},
		{dataType: types.Int32, value: uint64(42), expected: int32(42)},
		{dataType: types.Int32, value: int32(42), expected: int32(42)},
		{dataType: types.Int32, value: int64(42), expected: int32(42)},
		{dataType: types.Int32, value: float32(42), expected: int32(42)},
		{dataType: types.Int32, value: float64(42), expected: int32(42)},
		{dataType: types.Int64, value: uint32(42), expected: int64(42)},
		{dataType: types.Int64, value: uint64(42), expected: int64(42)},
		{dataType: types.Int64, value: int32(42), expected: int64(42)},
		{dataType: types.Int64, value: int64(42), expected: int64(42)},
		{dataType: types.Int64, value: float32(42), expected: int64(42)},
		{dataType: types.Int64, value: float64(42), expected: int64(42)},
		{dataType: types.Uint32, value: uint32(42), expected: uint32(42)},
		{dataType: types.Uint32, value: uint64(42), expected: uint32(42)},
		{dataType: types.Uint32, value: int32(42), expected: uint32(42)},
		{dataType: types.Uint32, value: int64(42), expected: uint32(42)},
		{dataType: types.Uint32, value: float32(42), expected: uint32(42)},
		{dataType: types.Uint32, value: float64(42), expected: uint32(42)},
		{dataType: types.Uint64, value: uint32(42), expected: uint64(42)},
		{dataType: types.Uint64, value: uint64(42), expected: uint64(42)},
		{dataType: types.Uint64, value: int64(42), expected: uint64(42)},
		{dataType: types.Uint64, value: int32(42), expected: uint64(42)},
		{dataType: types.Uint64, value: float32(42), expected: uint64(42)},
		{dataType: types.Uint64, value: float64(42), expected: uint64(42)},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%T to %s", test.value, test.dataType), func(t *testing.T) {
			var scalar = &Scalar{
				Type:    test.dataType,
				Default: test.value,
			}

			scalar.Clean()

			if !reflect.DeepEqual(scalar.Default, test.expected) {
				t.Errorf("unexpected default value %T(%+v), expected %T(%+v)", scalar.Default, scalar.Default, test.expected, test.expected)
			}
		})
	}
}

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
