package metadata

import (
	"testing"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
)

func TestNewManager(t *testing.T) {
	resource := "mock"
	params := specs.Header{
		"example": &specs.Property{
			Name:  "example",
			Path:  "example",
			Label: labels.Optional,
			Template: specs.Template{
				Scalar: &specs.Scalar{
					Type: types.String,
				},
			},
		},
	}

	ctx := logger.WithLogger(broker.NewBackground())
	manager := NewManager(ctx, resource, params)
	if manager == nil {
		t.Fatal("undefined manager")
	}
}

func TestManagerMarshal(t *testing.T) {
	resource := "mock"

	tests := map[string]func() (specs.Header, references.Store, MD){
		"simple": func() (specs.Header, references.Store, MD) {
			header := specs.Header{
				"example": &specs.Property{
					Name:  "example",
					Path:  "example",
					Label: labels.Optional,
					Template: specs.Template{
						Scalar: &specs.Scalar{
							Default: "hello",
							Type:    types.String,
						},
					},
				},
			}

			store := references.NewReferenceStore(1)

			expected := MD{
				"example": "hello",
			}

			return header, store, expected
		},
		"reference": func() (specs.Header, references.Store, MD) {
			header := specs.Header{
				"example": &specs.Property{
					Name:  "example",
					Path:  "example",
					Label: labels.Optional,
					Template: specs.Template{
						Scalar: &specs.Scalar{
							Type: types.String,
						},
						Reference: &specs.PropertyReference{
							Resource: "input",
							Path:     "value",
						},
					},
				},
			}

			store := references.NewReferenceStore(1)
			store.StoreValue("input", "value", "message")

			expected := MD{
				"example": "hello",
			}

			return header, store, expected
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			header, store, expected := test()
			ctx := logger.WithLogger(broker.NewBackground())
			manager := NewManager(ctx, resource, header)

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

	tests := map[string]func() (specs.Header, MD, MD){
		"simple": func() (specs.Header, MD, MD) {
			params := specs.Header{
				"example": &specs.Property{
					Name:  "example",
					Path:  "example",
					Label: labels.Optional,
					Template: specs.Template{
						Scalar: &specs.Scalar{
							Type: types.String,
						},
					},
				},
			}

			input := MD{
				"example": "hello",
			}

			expected := MD{
				"example": "hello",
			}

			return params, input, expected
		},
		"case insensitive": func() (specs.Header, MD, MD) {
			params := specs.Header{
				"Example": &specs.Property{
					Name:  "example",
					Path:  "example",
					Label: labels.Optional,
					Template: specs.Template{
						Scalar: &specs.Scalar{
							Type: types.String,
						},
					},
				},
			}

			input := MD{
				"example": "hello",
			}

			expected := MD{
				"example": "hello",
			}

			return params, input, expected
		},
		"unnessasery allocation": func() (specs.Header, MD, MD) {
			params := specs.Header{
				"example": &specs.Property{
					Name:  "example",
					Path:  "example",
					Label: labels.Optional,
					Template: specs.Template{
						Scalar: &specs.Scalar{
							Type: types.String,
						},
					},
				},
			}

			input := MD{
				"example": "hello",
				"unknown": "hello",
			}

			expected := MD{
				"example": "hello",
			}

			return params, input, expected
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			header, input, expected := test()
			ctx := logger.WithLogger(broker.NewBackground())
			manager := NewManager(ctx, resource, header)

			store := references.NewReferenceStore(len(input))
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

				_, has := expected[key]
				if !has {
					t.Fatalf("unexpected header key %s", key)
				}
			}
		})
	}
}
