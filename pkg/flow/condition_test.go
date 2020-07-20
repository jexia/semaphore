package flow

import (
	"errors"
	"testing"

	"github.com/jexia/maestro/pkg/conditions"
	"github.com/jexia/maestro/pkg/core/instance"
	"github.com/jexia/maestro/pkg/functions"
	"github.com/jexia/maestro/pkg/refs"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/labels"
	"github.com/jexia/maestro/pkg/specs/types"
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
					Fn: func(store refs.Store) error {
						store.StoreValue("", ".", 1)
						return nil
					},
					Returns: &specs.Property{
						Type:  types.String,
						Label: labels.Optional,
					},
				},
			},
			condition: "{{ stack.first:. }} == {{ input:id }}",
			expected:  true,
		},
	}

	store := refs.NewReferenceStore(5)
	store.StoreValue("input", "id", 1)
	store.StoreValue("input", "name", "john")
	store.StoreValue("input", "age", 99)
	store.StoreValue("input", "city", "Amsterdam")
	store.StoreValue("input", "country", "The Netherlands")

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := instance.NewContext()
			spec, err := conditions.NewEvaluableExpression(ctx, test.condition)
			if err != nil {
				t.Fatal(err)
			}

			condition := NewCondition(test.stack, spec)
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
					Fn: func(store refs.Store) error {
						store.StoreValue("", ".", 1)
						return errors.New("unexpected error")
					},
					Returns: &specs.Property{
						Type:  types.String,
						Label: labels.Optional,
					},
				},
			},
			condition: "{{ stack.first:. }} == {{ input:id }}",
		},
	}

	store := refs.NewReferenceStore(5)
	store.StoreValue("input", "id", 1)
	store.StoreValue("input", "name", "john")
	store.StoreValue("input", "age", 99)
	store.StoreValue("input", "city", "Amsterdam")
	store.StoreValue("input", "country", "The Netherlands")

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := instance.NewContext()
			spec, err := conditions.NewEvaluableExpression(ctx, test.condition)
			if err != nil {
				t.Fatal(err)
			}

			condition := NewCondition(test.stack, spec)
			_, err = condition.Eval(ctx, store)
			if err == nil {
				t.Fatal("unexpected pass")
			}
		})
	}
}
