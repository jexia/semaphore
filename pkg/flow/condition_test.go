package flow

import (
	"errors"
	"testing"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/conditions"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
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
						store.StoreValue("", ".", 1)
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

	store := references.NewReferenceStore(5)
	store.StoreValue("input", "id", 1)
	store.StoreValue("input", "name", "john")
	store.StoreValue("input", "age", 99)
	store.StoreValue("input", "city", "Amsterdam")
	store.StoreValue("input", "country", "The Netherlands")

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
						store.StoreValue("", ".", 1)
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

	store := references.NewReferenceStore(5)
	store.StoreValue("input", "id", 1)
	store.StoreValue("input", "name", "john")
	store.StoreValue("input", "age", 99)
	store.StoreValue("input", "city", "Amsterdam")
	store.StoreValue("input", "country", "The Netherlands")

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
