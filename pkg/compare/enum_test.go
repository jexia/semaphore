package compare

import (
	"testing"

	"github.com/jexia/semaphore/pkg/specs"
)

func TestCompareEnumValues(t *testing.T) {
	newValue := func() *specs.EnumValue {
		return &specs.EnumValue{
			Key:      "a",
			Position: 1,
		}
	}

	shouldBeEqual := func(t *testing.T, expected, given *specs.EnumValue) {
		got := CompareEnumValues(expected, given)
		if got != nil {
			t.Fatalf("CompareEnumValues() returns unexpected error: %v", got)
		}
	}

	shouldNotBeEqual := func(t *testing.T, expected, given *specs.EnumValue) {
		if CompareEnumValues(expected, given) == nil {
			t.Fatalf("CompareEnumValues() = nil, want an error")
		}
	}

	t.Run("should be equal", func(t *testing.T) {
		valA, valB := newValue(), newValue()

		shouldBeEqual(t, valA, valB)
	})

	t.Run("nils should be equal", func(t *testing.T) {
		shouldBeEqual(t, nil, nil)
	})

	t.Run("nil and not nil should not be equal", func(t *testing.T) {
		shouldNotBeEqual(t, newValue(), nil)
		shouldNotBeEqual(t, nil, newValue())
	})

	t.Run("values with different keys should not be equal", func(t *testing.T) {
		valA, valB := newValue(), newValue()

		valB.Key = "b"

		shouldNotBeEqual(t, valA, valB)
	})

	t.Run("values with different position should not be equal", func(t *testing.T) {
		valA, valB := newValue(), newValue()

		valB.Position = 42

		shouldNotBeEqual(t, valA, valB)
	})
}

func TestCompareEnums(t *testing.T) {
	createEnum := func() *specs.Enum {
		valA := &specs.EnumValue{
			Key:         "a",
			Position:    0,
			Description: "",
		}

		valB := &specs.EnumValue{
			Key:         "b",
			Position:    1,
			Description: "",
		}

		return &specs.Enum{
			Name:        "letters",
			Description: "",
			Keys:        map[string]*specs.EnumValue{valA.Key: valA, valB.Key: valB},
			Positions:   map[int32]*specs.EnumValue{valA.Position: valA, valB.Position: valB},
		}
	}

	shouldBeEqual := func(t *testing.T, expected, given *specs.Enum) {
		got := CompareEnums(expected, given)
		if got != nil {
			t.Fatalf("CompareEnums() returns unexpected error: %v", got)
		}
	}

	shouldNotBeEqual := func(t *testing.T, expected, given *specs.Enum) {
		if CompareEnums(expected, given) == nil {
			t.Fatalf("CompareEnums() = nil, want an error")
		}
	}

	t.Run("filled enums should be equal", func(t *testing.T) {
		enum, another := createEnum(), createEnum()

		shouldBeEqual(t, enum, another)
	})

	t.Run("empty enums should be equal", func(t *testing.T) {
		enum, another := createEnum(), createEnum()

		enum.Keys = nil
		enum.Positions = nil

		another.Keys = nil
		another.Positions = nil

		shouldBeEqual(t, enum, another)
	})

	t.Run("nil enums should be equal", func(t *testing.T) {
		shouldBeEqual(t, nil, nil)
	})

	t.Run("nil and not nil enums should not be equal", func(t *testing.T) {
		shouldNotBeEqual(t, nil, &specs.Enum{})
		shouldNotBeEqual(t, &specs.Enum{}, nil)
	})

	t.Run("enums with nil items under Keys/Positions should be equal", func(t *testing.T) {
		enum, another := createEnum(), createEnum()

		enum.Keys = map[string]*specs.EnumValue{"a": nil, "b": nil}
		enum.Positions = map[int32]*specs.EnumValue{0: nil, 1: nil}

		another.Keys = map[string]*specs.EnumValue{"a": nil, "b": nil}
		another.Positions = map[int32]*specs.EnumValue{0: nil, 1: nil}

		shouldBeEqual(t, enum, another)
	})

	t.Run("enums with different names should not be equal", func(t *testing.T) {
		enum, another := createEnum(), createEnum()

		another.Name = "breeds"

		shouldNotBeEqual(t, enum, another)
	})

	t.Run("enums with different size of Keys should not be equal", func(t *testing.T) {
		enum, another := createEnum(), createEnum()

		another.Keys["c"] = &specs.EnumValue{}

		shouldNotBeEqual(t, enum, another)
	})

	t.Run("enums with different size of Positions should not be equal", func(t *testing.T) {
		enum, another := createEnum(), createEnum()

		another.Positions[3] = &specs.EnumValue{}

		shouldNotBeEqual(t, enum, another)
	})

	t.Run("enums with different values under positions should not be equal", func(t *testing.T) {
		enum, another := createEnum(), createEnum()

		another.Positions[0] = &specs.EnumValue{}

		shouldNotBeEqual(t, enum, another)
	})

	t.Run("enums with different values under keys should not be equal", func(t *testing.T) {
		enum, another := createEnum(), createEnum()

		another.Keys["a"] = &specs.EnumValue{}

		shouldNotBeEqual(t, enum, another)
	})
}
