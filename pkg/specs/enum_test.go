package specs

import "testing"

func TestEnumValueCompare(t *testing.T) {
	newValue := func() *EnumValue {
		return &EnumValue{
			Key:      "a",
			Position: 1,
		}
	}

	shouldBeEqual := func(t *testing.T, expected, given *EnumValue) {
		if err := given.Compare(expected); err != nil {
			t.Fatalf("CompareEnumValues() returns unexpected error: %v", err)
		}
	}

	shouldNotBeEqual := func(t *testing.T, expected, given *EnumValue) {
		if given.Compare(expected) == nil {
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

func TestEnumCompare(t *testing.T) {
	createEnum := func() *Enum {
		valA := &EnumValue{
			Key:         "a",
			Position:    0,
			Description: "",
		}

		valB := &EnumValue{
			Key:         "b",
			Position:    1,
			Description: "",
		}

		return &Enum{
			Name:        "letters",
			Description: "",
			Keys:        map[string]*EnumValue{valA.Key: valA, valB.Key: valB},
			Positions:   map[int32]*EnumValue{valA.Position: valA, valB.Position: valB},
		}
	}

	shouldBeEqual := func(t *testing.T, expected, given *Enum) {
		got := given.Compare(expected)
		if got != nil {
			t.Fatalf("CompareEnums() returns unexpected error: %v", got)
		}
	}

	shouldNotBeEqual := func(t *testing.T, expected, given *Enum) {
		if given.Compare(expected) == nil {
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
		shouldNotBeEqual(t, nil, &Enum{})
		shouldNotBeEqual(t, &Enum{}, nil)
	})

	t.Run("enums with nil items under Keys/Positions should be equal", func(t *testing.T) {
		enum, another := createEnum(), createEnum()

		enum.Keys = map[string]*EnumValue{"a": nil, "b": nil}
		enum.Positions = map[int32]*EnumValue{0: nil, 1: nil}

		another.Keys = map[string]*EnumValue{"a": nil, "b": nil}
		another.Positions = map[int32]*EnumValue{0: nil, 1: nil}

		shouldBeEqual(t, enum, another)
	})

	t.Run("enums with different names should not be equal", func(t *testing.T) {
		enum, another := createEnum(), createEnum()

		another.Name = "breeds"

		shouldNotBeEqual(t, enum, another)
	})

	t.Run("enums with different size of Keys should not be equal", func(t *testing.T) {
		enum, another := createEnum(), createEnum()

		another.Keys["c"] = &EnumValue{}

		shouldNotBeEqual(t, enum, another)
	})

	t.Run("enums with different size of Positions should not be equal", func(t *testing.T) {
		enum, another := createEnum(), createEnum()

		another.Positions[3] = &EnumValue{}

		shouldNotBeEqual(t, enum, another)
	})

	t.Run("enums with different values under positions should not be equal", func(t *testing.T) {
		enum, another := createEnum(), createEnum()

		another.Positions[0] = &EnumValue{}

		shouldNotBeEqual(t, enum, another)
	})

	t.Run("enums with different values under keys should not be equal", func(t *testing.T) {
		enum, another := createEnum(), createEnum()

		another.Keys["a"] = &EnumValue{}

		shouldNotBeEqual(t, enum, another)
	})
}
