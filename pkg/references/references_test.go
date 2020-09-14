package references

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/compare"
	"github.com/jexia/semaphore/pkg/prettyerr"
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
				stack, perr := prettyerr.Prettify(err)
				if perr != nil {
					t.Fatal(perr)
				}

				if _, err := prettyerr.TextFormatter(stack, prettyerr.DefaultTextFormat); err != nil {
					t.Fatal(err)
				}

				if err.Error() != collection.Exception.Message {
					t.Fatalf("unexpected error message %s, expected %s", err.Error(), collection.Exception.Message)
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
				stack, perr := prettyerr.Prettify(err)
				if perr != nil {
					t.Fatal(perr)
				}

				if _, err := prettyerr.TextFormatter(stack, prettyerr.DefaultTextFormat); err != nil {
					t.Fatal(err)
				}

				if err.Error() != collection.Exception.Message {
					t.Fatalf("unexpected error message %s, expected %s", err.Error(), collection.Exception.Message)
				}
			}
		})
	}
}

func TestScopeNestedReferencesNil(t *testing.T) {
	ScopeNestedReferences(nil, nil)
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
				Repeated: []*specs.Property{
					{
						Name: "key",
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
				Repeated: []*specs.Property{
					{
						Name: "key",
						Path: "key",
						Repeated: []*specs.Property{
							{
								Name: "nested",
								Path: "key.nested",
							},
						},
					},
				},
			},
			target: &specs.Property{
				Repeated: []*specs.Property{
					{
						Name:      "key",
						Path:      "key",
						Reference: reference,
					},
				},
			},
		},
		"partial": {
			source: &specs.Property{
				Repeated: []*specs.Property{
					{
						Name: "key",
						Path: "key",
						Repeated: []*specs.Property{
							{
								Name: "first",
								Path: "key.first",
							},
							{
								Name: "second",
								Path: "key.second",
							},
						},
					},
				},
			},
			target: &specs.Property{
				Repeated: []*specs.Property{
					{
						Name:      "key",
						Path:      "key",
						Reference: reference,
						Repeated: []*specs.Property{
							{
								Name:      "second",
								Path:      "key.second",
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
				if len(target.Repeated) != len(source.Repeated) {
					t.Fatalf("unexpected length %d (%+v), expected %d (%s)(%+v).", len(target.Repeated), target.Repeated, len(source.Repeated), source.Path, source.Repeated)
				}

				for _, item := range source.Repeated {
					target := target.Repeated.Get(item.Name)
					if target == nil {
						t.Fatalf("target does not have nested key %s", item.Name)
					}

					lookup(item, target)
				}
			}

			lookup(test.source, test.target)
		})
	}
}
