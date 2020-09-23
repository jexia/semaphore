package conditions

import (
	"testing"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/specs"
)

func TestNewEvaluableExpression(t *testing.T) {
	type test struct {
		raw    string
		params map[string]*specs.Property
	}

	tests := []test{
		{
			raw: "{{ input:id }} == {{ input:id }}",
			params: map[string]*specs.Property{
				"input:id": {
					Reference: &specs.PropertyReference{
						Resource: "input",
						Path:     "id",
					},
				},
			},
		},
		{
			raw: "({{ input:id }} == {{ input:id }}) || {{ input:name }}",
			params: map[string]*specs.Property{
				"input:id": {
					Reference: &specs.PropertyReference{
						Resource: "input",
						Path:     "id",
					},
				},
				"input:name": {
					Reference: &specs.PropertyReference{
						Resource: "input",
						Path:     "name",
					},
				},
			},
		},
		{
			raw: "({{ resource:id }} == {{ input:id }})",
			params: map[string]*specs.Property{
				"input:id": {
					Reference: &specs.PropertyReference{
						Resource: "input",
						Path:     "id",
					},
				},
				"resource:id": {
					Reference: &specs.PropertyReference{
						Resource: "resource",
						Path:     "id",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.raw, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())
			condition, err := NewEvaluableExpression(ctx, test.raw)
			if err != nil {
				t.Fatal(err)
			}

			for key, param := range condition.Params.Params {
				expected, has := test.params[key]
				if !has {
					t.Fatalf("unexpected result, expected %s to be set", key)
				}

				if expected.Reference != nil && param.Reference == nil {
					t.Fatalf("unexpected reference %s, reference not set", key)
				}

				if expected.Reference != nil {
					if param.Reference.Resource != expected.Reference.Resource {
						t.Fatalf("unexpected resource '%+v', expected '%+v'", param.Reference.Resource, expected.Reference.Resource)
					}

					if param.Reference.Path != expected.Reference.Path {
						t.Fatalf("unexpected path '%+v', expected '%+v'", param.Reference.Path, expected.Reference.Path)
					}
				}

				if param.Type() != expected.Type() {
					t.Fatalf("unexpected type '%+v', expected '%+v'", param.Type(), expected.Type())
				}

				if param.Label != expected.Label {
					t.Fatalf("unexpected label '%+v', expected '%+v'", param.Label, expected.Label)
				}
			}
		})
	}
}

func TestInvalidExpressions(t *testing.T) {
	tests := []string{
		"( {{ input:id }}",
		"== {{ input:id }}",
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())
			_, err := NewEvaluableExpression(ctx, test)
			if err == nil {
				t.Fatal("unexpected pass")
			}
		})
	}
}

func TestInvalidReference(t *testing.T) {
	tests := []string{
		"{{ input:id.. }}",
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())
			_, err := NewEvaluableExpression(ctx, test)
			if err == nil {
				t.Fatal("unexpected pass")
			}
		})
	}
}
