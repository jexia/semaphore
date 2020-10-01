package specs

import (
	"testing"

	"github.com/jexia/semaphore/pkg/specs/metadata"
)

func TestPropertyReferenceClone(t *testing.T) {
	reference := &PropertyReference{
		Meta:     metadata.WithValue(nil, nil, nil),
		Resource: "resource",
		Path:     "path",
		Property: &Property{},
	}

	result := reference.Clone()
	if result == nil {
		t.Error("unexpected result, expected property reference clone to be returned")
	}

	if result.Meta != reference.Meta {
		t.Errorf("unexpected meta %+v, expected %+v", result.Meta, reference.Meta)
	}

	if result.Resource != reference.Resource {
		t.Errorf("unexpected resource %+v", result.Resource)
	}

	if result.Path != reference.Path {
		t.Errorf("unexpected path %+v", result.Path)
	}

	if result.Property != nil {
		t.Errorf("unexpected property %+v", result.Property)
	}
}

func TestPropertyReferenceCloneNilValue(t *testing.T) {
	var reference *PropertyReference
	result := reference.Clone()
	if result != nil {
		t.Errorf("unexpected result %+v", result)
	}
}

func TestPropertyReferenceString(t *testing.T) {
	t.Parallel()

	tests := map[string]*PropertyReference{
		"resource:path": {
			Resource: "resource",
			Path:     "path",
		},
		"resource:nested.path": {
			Resource: "resource",
			Path:     "nested.path",
		},
		"resource.prop:path": {
			Resource: "resource.prop",
			Path:     "path",
		},
		":": {},
		"":  nil,
	}

	for expected, reference := range tests {
		t.Run(expected, func(t *testing.T) {
			result := reference.String()
			if result != expected {
				t.Fatalf("unexpected result %s, expected %s", result, expected)
			}
		})
	}
}
