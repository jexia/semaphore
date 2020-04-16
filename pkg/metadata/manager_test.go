package metadata

import (
	"testing"

	"github.com/jexia/maestro/pkg/refs"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/labels"
	"github.com/jexia/maestro/pkg/specs/types"
)

func TestNewManager(t *testing.T) {
	resource := "mock"
	params := specs.Header{
		"example": &specs.Property{
			Name:  "example",
			Path:  "example",
			Type:  types.String,
			Label: labels.Optional,
		},
	}

	manager := NewManager(resource, params)
	if manager == nil {
		t.Fatal("undefined manager")
	}
}

func TestManagerMarshal(t *testing.T) {
	resource := "mock"

	tests := map[string]func() (specs.Header, refs.Store, MD){
		"simple": func() (specs.Header, refs.Store, MD) {
			header := specs.Header{
				"example": &specs.Property{
					Name:    "example",
					Path:    "example",
					Default: "hello",
					Type:    types.String,
					Label:   labels.Optional,
				},
			}

			store := refs.NewReferenceStore(1)

			expected := MD{
				"example": "hello",
			}

			return header, store, expected
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			header, store, expected := test()
			manager := NewManager(resource, header)

			result := manager.Marshal(store)

			for key, val := range result {
				expected, has := expected[key]
				if !has {
					t.Errorf("unexpected key %s in result", key)
				}

				if val != val {
					t.Errorf("unexpected value %s, expected %s", val, expected)
				}
			}
		})
	}
}

func TestManagerUnmarshal(t *testing.T) {
	resource := "mock"

	tests := map[string]func() (specs.Header, MD){
		"simple": func() (specs.Header, MD) {
			params := specs.Header{
				"example": &specs.Property{
					Name:  "example",
					Path:  "example",
					Type:  types.String,
					Label: labels.Optional,
				},
			}

			input := MD{
				"example": "hello",
			}

			return params, input
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			header, input := test()
			manager := NewManager(resource, header)

			store := refs.NewReferenceStore(len(input))
			manager.Unmarshal(input, store)

			for key, prop := range header {
				if prop.Reference == nil {
					continue
				}

				ref := store.Load(resource, key)
				if ref == nil {
					t.Fatalf("reference not set %s", key)
				}

				str, is := ref.Value.(string)
				if !is {
					t.Fatalf("reference value is not a string %+v", ref.Value)
				}

				if str != input[key] {
					t.Fatalf("unexpected value %s, expected %s", str, input[key])
				}
			}
		})
	}
}
