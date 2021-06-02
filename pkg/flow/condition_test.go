package flow

import (
	"errors"
	"testing"

	"github.com/jexia/semaphore/v2/pkg/broker"
	"github.com/jexia/semaphore/v2/pkg/broker/logger"
	"github.com/jexia/semaphore/v2/pkg/conditions"
	"github.com/jexia/semaphore/v2/pkg/functions"
	"github.com/jexia/semaphore/v2/pkg/references"
	"github.com/jexia/semaphore/v2/pkg/specs"
	"github.com/jexia/semaphore/v2/pkg/specs/labels"
	"github.com/jexia/semaphore/v2/pkg/specs/types"
)

func TestConditionEvaluation(t *testing.T) {
	type test struct {
		stack     functions.Stack
		condition string
		expected  bool
	}

	tests := map[string]test{
		"simple": {
			stack:     nil,
			condition: "{{ input:id }} == {{ input:id }}",
			expected:  true,
		},
		"not": {
			stack:     nil,
			condition: "{{ input:id }} != {{ input:name }}",
			expected:  true,
		},
		"false": {
			stack:     nil,
			condition: "{{ input:id }} == {{ input:age }}",
			expected:  false,
		},
		"property": {
			stack:     nil,
			condition: "{{ input:id }}",
			expected:  true,
		},
		"functions": {
			stack: functions.Stack{
				"first": &functions.Function{
					Fn: func(store references.Store) error {
						store.Store("", &references.Reference{Value: 1})
						return nil
					},
					Returns: &specs.Property{
						Label: labels.Optional,
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Type: types.String,
							},
						},
					},
				},
			},
			condition: "{{ stack.first:. }} == {{ input:id }}",
			expected:  true,
		},
	}

	store := references.NewStore(5)
	store.Define(specs.ResourcePath("input"), 5)
	store.Store(specs.ResourcePath("input", "id"), &references.Reference{Value: 1})
	store.Store(specs.ResourcePath("input", "name"), &references.Reference{Value: "john"})
	store.Store(specs.ResourcePath("input", "age"), &references.Reference{Value: 99})
	store.Store(specs.ResourcePath("input", "city"), &references.Reference{Value: "Amsterdam"})
	store.Store(specs.ResourcePath("input", "country"), &references.Reference{Value: "The Netherlands"})

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())
			exp, err := conditions.NewEvaluableExpression(ctx, test.condition)
			if err != nil {
				t.Fatal(err)
			}

			condition := NewCondition(test.stack, exp)
			result, err := condition.Eval(ctx, store)
			if err != nil {
				t.Fatal(err)
			}

			if result != test.expected {
				t.Fatalf("unexpected result %t, expected %t", result, test.expected)
			}
		})
	}
}

func TestInvalidConditionEvaluation(t *testing.T) {
	type test struct {
		stack     functions.Stack
		condition string
	}

	tests := map[string]test{
		"invalid types": {
			stack:     nil,
			condition: "{{ input:name }} > {{ input:id }}",
		},
		"function error": {
			stack: functions.Stack{
				"first": &functions.Function{
					Fn: func(store references.Store) error {
						store.Store("", &references.Reference{Value: 1})
						return errors.New("unexpected error")
					},
					Returns: &specs.Property{
						Label: labels.Optional,
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Type: types.String,
							},
						},
					},
				},
			},
			condition: "{{ stack.first:. }} == {{ input:id }}",
		},
	}

	store := references.NewStore(5)
	store.Define(specs.ResourcePath("input"), 5)
	store.Store(specs.ResourcePath("input", "id"), &references.Reference{Value: 1})
	store.Store(specs.ResourcePath("input", "name"), &references.Reference{Value: "john"})
	store.Store(specs.ResourcePath("input", "age"), &references.Reference{Value: 99})
	store.Store(specs.ResourcePath("input", "city"), &references.Reference{Value: "Amsterdam"})
	store.Store(specs.ResourcePath("input", "country"), &references.Reference{Value: "The Netherlands"})

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())
			exp, err := conditions.NewEvaluableExpression(ctx, test.condition)
			if err != nil {
				t.Fatal(err)
			}

			condition := NewCondition(test.stack, exp)
			_, err = condition.Eval(ctx, store)
			if err == nil {
				t.Fatal("unexpected pass")
			}
		})
	}
}
