package references

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/compare"
	"github.com/jexia/semaphore/pkg/providers"
	"github.com/jexia/semaphore/pkg/providers/hcl"
	"github.com/jexia/semaphore/pkg/providers/mock"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

const (
	pass = "pass"
	fail = "fail"
)

func TestUnmarshalFile(t *testing.T) {
	path, err := filepath.Abs("./tests/*.hcl")
	if err != nil {
		t.Fatal(err)
	}

	ctx := logger.WithLogger(broker.NewBackground())
	files, err := providers.ResolvePath(ctx, []string{}, path)
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		t.Run(file.Name(), func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())

			flows, err := hcl.FlowsResolver(file.Path)(ctx)
			if err != nil {
				t.Fatal(err)
			}

			clean := file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]
			path := filepath.Join(filepath.Dir(file.Path), clean+".yaml")

			collection, err := mock.CollectionResolver(path)
			if err != nil {
				t.Fatal(err)
			}

			services, err := mock.ServicesResolver(path)(ctx)
			if err != nil {
				t.Fatal(err)
			}

			schemas, err := mock.SchemaResolver(path)(ctx)
			if err != nil {
				t.Fatal(err)
			}

			err = func() (err error) {
				err = providers.ResolveSchemas(ctx, services, schemas, flows)
				if err != nil {
					return err
				}

				err = Resolve(ctx, flows)
				if err != nil {
					return err
				}

				return nil
			}()

			if strings.HasSuffix(clean, pass) && err != nil {
				t.Fatalf("expected test to pass but failed instead %s, %v", file.Name(), err)
			}

			if strings.HasSuffix(clean, fail) && err != nil {
				if err.Error() != collection.Exception.Message {
					t.Fatalf("unexpected error message %s, expected %s", err, collection.Exception.Message)
				}

				return
			}

			err = compare.Types(ctx, services, schemas, flows)

			if strings.HasSuffix(clean, pass) && err != nil {
				t.Fatalf("expected test to pass but failed instead %s, %v", file.Name(), err)
			}

			if strings.HasSuffix(clean, fail) && err == nil {
				t.Fatalf("expected test to fail but passed instead %s", file.Name())
			}

			if strings.HasSuffix(clean, fail) {
				if err.Error() != collection.Exception.Message {
					t.Fatalf("unexpected error message %s, expected %s", err, collection.Exception.Message)
				}
			}
		})
	}
}

func TestScopeNestedReferences(t *testing.T) {
	t.Parallel()

	type test struct {
		source *specs.Property
		target *specs.Property
	}

	reference := &specs.PropertyReference{
		Resource: "input",
	}

	tests := map[string]test{
		"root": {
			source: &specs.Property{
				Nested: map[string]*specs.Property{
					"key": {
						Path: "key",
						Type: types.String,
					},
				},
			},
			target: &specs.Property{
				Reference: reference,
			},
		},
		"nested": {
			source: &specs.Property{
				Nested: map[string]*specs.Property{
					"key": {
						Path: "key",
						Nested: map[string]*specs.Property{
							"nested": {
								Path: "key.nested",
							},
						},
					},
				},
			},
			target: &specs.Property{
				Nested: map[string]*specs.Property{
					"key": {
						Reference: reference,
					},
				},
			},
		},
		"partial": {
			source: &specs.Property{
				Nested: map[string]*specs.Property{
					"key": {
						Path: "key",
						Nested: map[string]*specs.Property{
							"first": {
								Path: "key.first",
							},
							"second": {
								Path: "key.second",
							},
						},
					},
				},
			},
			target: &specs.Property{
				Nested: map[string]*specs.Property{
					"key": {
						Reference: reference,
						Nested: map[string]*specs.Property{
							"second": {
								Reference: reference,
							},
						},
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ScopeNestedReferences(test.source, test.target)

			var lookup func(source *specs.Property, target *specs.Property)
			lookup = func(source *specs.Property, target *specs.Property) {
				if len(target.Nested) != len(source.Nested) {
					t.Fatalf("unexpected length %d (%+v), expected %d (%s)(%+v).", len(target.Nested), target.Nested, len(source.Nested), source.Path, source.Nested)
				}

				for key := range source.Nested {
					_, has := target.Nested[key]
					if !has {
						t.Fatalf("target does not have nested key %s", key)
					}

					lookup(source.Nested[key], target.Nested[key])
				}
			}

			lookup(test.source, test.target)
		})
	}
}
