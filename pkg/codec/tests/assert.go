package tests

import (
	"testing"

	"github.com/jexia/semaphore/v2/pkg/references"
	"github.com/jexia/semaphore/v2/pkg/specs/template"
)

// Expect contains expected value.
type Expect struct {
	Scalar interface{}
	Enum   *int32
}

// Assert checks if value under given path matches expected.
func Assert(t *testing.T, resource, path string, store references.Store, input Expect) {
	reference := store.Load(template.ResourcePath(resource, path))

	switch {
	case input.Enum != nil:
		if reference == nil {
			t.Fatalf("reference %q was expected to be set", path)
		}

		if reference.Enum == nil {
			t.Fatalf("reference %q was expected to have a enum value", path)
		}

		if *input.Enum != *reference.Enum {
			t.Errorf("reference %q was expected to have enum value [%d], not [%d]", path, *input.Enum, *reference.Enum)
		}
	case input.Scalar != nil:
		if reference == nil {
			t.Fatalf("reference %q was expected to be set", path)
		}

		if reference.Value != input.Scalar {
			t.Errorf("reference %q was expected to be %T(%v), got %T(%v)", path, input.Scalar, input.Scalar, reference.Value, reference.Value)
		}
	}
}
