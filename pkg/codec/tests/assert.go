package tests

import (
	"testing"

	"github.com/jexia/semaphore/pkg/references"
)

type Expect struct {
	Value    interface{}
	Enum     *int32
	Repeated []Expect
	Nested   map[string]Expect
}

func buildPath(prefix, property string) string {
	if prefix == "" {
		return property
	}

	return prefix + "." + property
}

func Assert(t *testing.T, resource, path string, store references.Store, input Expect) {
	ref := store.Load(resource, path)

	switch {
	case input.Nested != nil:
		for key, value := range input.Nested {
			Assert(t, resource, buildPath(path, key), store, value)
		}
	case input.Enum != nil:
		if ref == nil {
			t.Fatalf("reference %q was expected to be set", path)
		}

		if ref.Enum == nil {
			t.Fatalf("reference %q was expected to have a enum value", path)
		}

		if *input.Enum != *ref.Enum {
			t.Errorf("reference %q was expected to have enum value [%d], not [%d]", path, *input.Enum, *ref.Enum)
		}
	case input.Value != nil:
		if ref == nil {
			t.Fatalf("reference %q was expected to be set", path)
		}

		if ref.Value != input.Value {
			t.Errorf("reference %q was expected to be %T(%v), got %T(%v)", path, input.Value, input.Value, ref.Value, ref.Value)
		}
	case input.Repeated != nil:
		if ref == nil {
			t.Fatalf("reference %q was expected to be set", path)
		}

		if ref.Repeated == nil {
			t.Fatalf("reference %q was expected to have a repeated value", path)
		}

		if expected, actual := len(input.Repeated), len(ref.Repeated); actual != expected {
			t.Fatalf("invalid number of repeated values, expected %d, got %d", expected, actual)
		}

		for index, expected := range input.Repeated {
			Assert(t, "", "", ref.Repeated[index], expected)
		}
	}
}
