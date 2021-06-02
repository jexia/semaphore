package specs

import (
	"testing"

	"github.com/jexia/semaphore/v2/pkg/specs/metadata"
)

func TestParsePropertyReference(t *testing.T) {
	tests := map[string]Template{
		// Simple/default
		"input:message": {
			Reference: &PropertyReference{
				Resource: "input",
				Path:     "message",
			},
		},
		// Lowercase path
		"input.header:Authorization": {
			Reference: &PropertyReference{
				Resource: "input.header",
				Path:     "authorization",
			},
		},
		// No path
		"input:": {
			Reference: &PropertyReference{
				Resource: "input",
			},
		},
		// No path complex(er)
		"input.header:": {
			Reference: &PropertyReference{
				Resource: "input.header",
			},
		},
		// No delimiter
		"input": {
			Reference: &PropertyReference{
				Resource: "input",
			},
		},
	}

	for input, expected := range tests {
		t.Run(input, func(t *testing.T) {
			propertyReference := ParsePropertyReference(input)

			CompareTemplates(t, Template{
				Reference: propertyReference,
			}, expected)
		})
	}
}

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
