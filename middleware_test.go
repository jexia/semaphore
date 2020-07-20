package semaphore

import (
	"context"
	"testing"

	"github.com/jexia/semaphore/pkg/core/api"
	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/flow"
	"github.com/jexia/semaphore/pkg/refs"
	"github.com/jexia/semaphore/pkg/specs"
)

func TestWithMiddleware(t *testing.T) {
	middleware := func(instance.Context) ([]api.Option, error) {
		return nil, nil
	}

	client, err := New(WithMiddleware(middleware), WithMiddleware(middleware))
	if err != nil {
		t.Fatal(err)
	}

	if client.Options.Middleware == nil {
		t.Fatal("middleware not set")
	}

	if len(client.Options.Middleware) != 2 {
		t.Fatalf("unexpected middleware %d, expected 2", len(client.Options.Middleware))
	}
}

func TestAfterConstructorOption(t *testing.T) {
	fn := func(i *int) api.AfterConstructorHandler {
		return func(next api.AfterConstructor) api.AfterConstructor {
			return func(ctx instance.Context, flow *specs.Collection) error {
				*i++
				return next(ctx, flow)
			}
		}
	}

	type test struct {
		expected  int
		arguments func() (*int, []api.Option)
	}

	tests := map[string]test{
		"single": {
			expected: 1,
			arguments: func() (*int, []api.Option) {
				result := 0
				arguments := NewCollection(AfterConstructor(fn(&result)))

				return &result, arguments
			},
		},
		"multiple": {
			expected: 2,
			arguments: func() (*int, []api.Option) {
				result := 0
				arguments := NewCollection(AfterConstructor(fn(&result)), AfterConstructor(fn(&result)))

				return &result, arguments
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, options := test.arguments()
			client, err := New(options...)
			if err != nil {
				t.Fatal(err)
			}

			if client.Options.AfterConstructor == nil {
				t.Fatal("unexpected result expected option to be set")
			}

			if *result != test.expected {
				t.Fatalf("unexpected result %d, expected %d", *result, test.expected)
			}
		})
	}
}

func TestBeforeManagerDoOption(t *testing.T) {
	fn := func(i *int) flow.BeforeManagerHandler {
		return func(next flow.BeforeManager) flow.BeforeManager {
			return func(ctx context.Context, manager *flow.Manager, store refs.Store) (context.Context, error) {
				*i++
				return next(ctx, manager, store)
			}
		}
	}

	type test struct {
		expected  int
		arguments func() (*int, []api.Option)
	}

	tests := map[string]test{
		"single": {
			expected: 1,
			arguments: func() (*int, []api.Option) {
				result := 0
				arguments := NewCollection(BeforeManagerDo(fn(&result)))

				return &result, arguments
			},
		},
		"multiple": {
			expected: 2,
			arguments: func() (*int, []api.Option) {
				result := 0
				arguments := NewCollection(BeforeManagerDo(fn(&result)), BeforeManagerDo(fn(&result)))

				return &result, arguments
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, options := test.arguments()
			client, err := New(options...)
			if err != nil {
				t.Fatal(err)
			}

			if client.Options.BeforeManagerDo == nil {
				t.Fatal("unexpected result expected option to be set")
			}

			_, err = client.Options.BeforeManagerDo(nil, nil, nil)
			if err != nil {
				t.Fatal(err)
			}

			if *result != test.expected {
				t.Fatalf("unexpected result %d, expected %d", *result, test.expected)
			}
		})
	}
}

func TestBeforeManagerRollbackOption(t *testing.T) {
	fn := func(i *int) flow.BeforeManagerHandler {
		return func(next flow.BeforeManager) flow.BeforeManager {
			return func(ctx context.Context, manager *flow.Manager, store refs.Store) (context.Context, error) {
				*i++
				return next(ctx, manager, store)
			}
		}
	}

	type test struct {
		expected  int
		arguments func() (*int, []api.Option)
	}

	tests := map[string]test{
		"single": {
			expected: 1,
			arguments: func() (*int, []api.Option) {
				result := 0
				arguments := NewCollection(BeforeManagerRollback(fn(&result)))

				return &result, arguments
			},
		},
		"multiple": {
			expected: 2,
			arguments: func() (*int, []api.Option) {
				result := 0
				arguments := NewCollection(BeforeManagerRollback(fn(&result)), BeforeManagerRollback(fn(&result)))

				return &result, arguments
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, options := test.arguments()
			client, err := New(options...)
			if err != nil {
				t.Fatal(err)
			}

			if client.Options.BeforeManagerRollback == nil {
				t.Fatal("unexpected result expected option to be set")
			}

			_, err = client.Options.BeforeManagerRollback(nil, nil, nil)
			if err != nil {
				t.Fatal(err)
			}

			if *result != test.expected {
				t.Fatalf("unexpected result %d, expected %d", *result, test.expected)
			}
		})
	}
}

func TestAfterManagerRollbackOption(t *testing.T) {
	fn := func(i *int) flow.AfterManagerHandler {
		return func(next flow.AfterManager) flow.AfterManager {
			return func(ctx context.Context, manager *flow.Manager, store refs.Store) (context.Context, error) {
				*i++
				return next(ctx, manager, store)
			}
		}
	}

	type test struct {
		expected  int
		arguments func() (*int, []api.Option)
	}

	tests := map[string]test{
		"single": {
			expected: 1,
			arguments: func() (*int, []api.Option) {
				result := 0
				arguments := NewCollection(AfterManagerRollback(fn(&result)))

				return &result, arguments
			},
		},
		"multiple": {
			expected: 2,
			arguments: func() (*int, []api.Option) {
				result := 0
				arguments := NewCollection(AfterManagerRollback(fn(&result)), AfterManagerRollback(fn(&result)))

				return &result, arguments
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, options := test.arguments()
			client, err := New(options...)
			if err != nil {
				t.Fatal(err)
			}

			if client.Options.AfterManagerRollback == nil {
				t.Fatal("unexpected result expected option to be set")
			}

			_, err = client.Options.AfterManagerRollback(nil, nil, nil)
			if err != nil {
				t.Fatal(err)
			}

			if *result != test.expected {
				t.Fatalf("unexpected result %d, expected %d", *result, test.expected)
			}
		})
	}
}

func TestAfterManagerDoOption(t *testing.T) {
	fn := func(i *int) flow.AfterManagerHandler {
		return func(next flow.AfterManager) flow.AfterManager {
			return func(ctx context.Context, manager *flow.Manager, store refs.Store) (context.Context, error) {
				*i++
				return next(ctx, manager, store)
			}
		}
	}

	type test struct {
		expected  int
		arguments func() (*int, []api.Option)
	}

	tests := map[string]test{
		"single": {
			expected: 1,
			arguments: func() (*int, []api.Option) {
				result := 0
				arguments := NewCollection(AfterManagerDo(fn(&result)))

				return &result, arguments
			},
		},
		"multiple": {
			expected: 2,
			arguments: func() (*int, []api.Option) {
				result := 0
				arguments := NewCollection(AfterManagerDo(fn(&result)), AfterManagerDo(fn(&result)))

				return &result, arguments
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, options := test.arguments()
			client, err := New(options...)
			if err != nil {
				t.Fatal(err)
			}

			if client.Options.AfterManagerDo == nil {
				t.Fatal("unexpected result expected option to be set")
			}

			_, err = client.Options.AfterManagerDo(nil, nil, nil)
			if err != nil {
				t.Fatal(err)
			}

			if *result != test.expected {
				t.Fatalf("unexpected result %d, expected %d", *result, test.expected)
			}
		})
	}
}

func TestBeforeNodeDoOption(t *testing.T) {
	fn := func(i *int) flow.BeforeNodeHandler {
		return func(next flow.BeforeNode) flow.BeforeNode {
			return func(ctx context.Context, node *flow.Node, tracker *flow.Tracker, processes *flow.Processes, store refs.Store) (context.Context, error) {
				*i++
				return next(ctx, node, tracker, processes, store)
			}
		}
	}

	type test struct {
		expected  int
		arguments func() (*int, []api.Option)
	}

	tests := map[string]test{
		"single": {
			expected: 1,
			arguments: func() (*int, []api.Option) {
				result := 0
				arguments := NewCollection(BeforeNodeDo(fn(&result)))

				return &result, arguments
			},
		},
		"multiple": {
			expected: 2,
			arguments: func() (*int, []api.Option) {
				result := 0
				arguments := NewCollection(BeforeNodeDo(fn(&result)), BeforeNodeDo(fn(&result)))

				return &result, arguments
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, options := test.arguments()
			client, err := New(options...)
			if err != nil {
				t.Fatal(err)
			}

			if client.Options.BeforeNodeDo == nil {
				t.Fatal("unexpected result expected option to be set")
			}

			_, err = client.Options.BeforeNodeDo(nil, nil, nil, nil, nil)
			if err != nil {
				t.Fatal(err)
			}

			if *result != test.expected {
				t.Fatalf("unexpected result %d, expected %d", *result, test.expected)
			}
		})
	}
}

func TestBeforeNodeRollbackOption(t *testing.T) {
	fn := func(i *int) flow.BeforeNodeHandler {
		return func(next flow.BeforeNode) flow.BeforeNode {
			return func(ctx context.Context, node *flow.Node, tracker *flow.Tracker, processes *flow.Processes, store refs.Store) (context.Context, error) {
				*i++
				return next(ctx, node, tracker, processes, store)
			}
		}
	}

	type test struct {
		expected  int
		arguments func() (*int, []api.Option)
	}

	tests := map[string]test{
		"single": {
			expected: 1,
			arguments: func() (*int, []api.Option) {
				result := 0
				arguments := NewCollection(BeforeNodeRollback(fn(&result)))

				return &result, arguments
			},
		},
		"multiple": {
			expected: 2,
			arguments: func() (*int, []api.Option) {
				result := 0
				arguments := NewCollection(BeforeNodeRollback(fn(&result)), BeforeNodeRollback(fn(&result)))

				return &result, arguments
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, options := test.arguments()
			client, err := New(options...)
			if err != nil {
				t.Fatal(err)
			}

			if client.Options.BeforeNodeRollback == nil {
				t.Fatal("unexpected result expected option to be set")
			}

			_, err = client.Options.BeforeNodeRollback(nil, nil, nil, nil, nil)
			if err != nil {
				t.Fatal(err)
			}

			if *result != test.expected {
				t.Fatalf("unexpected result %d, expected %d", *result, test.expected)
			}
		})
	}
}

func TestAfterNodeRollbackOption(t *testing.T) {
	fn := func(i *int) flow.AfterNodeHandler {
		return func(next flow.AfterNode) flow.AfterNode {
			return func(ctx context.Context, node *flow.Node, tracker *flow.Tracker, processes *flow.Processes, store refs.Store) (context.Context, error) {
				*i++
				return next(ctx, node, tracker, processes, store)
			}
		}
	}

	type test struct {
		expected  int
		arguments func() (*int, []api.Option)
	}

	tests := map[string]test{
		"single": {
			expected: 1,
			arguments: func() (*int, []api.Option) {
				result := 0
				arguments := NewCollection(AfterNodeRollback(fn(&result)))

				return &result, arguments
			},
		},
		"multiple": {
			expected: 2,
			arguments: func() (*int, []api.Option) {
				result := 0
				arguments := NewCollection(AfterNodeRollback(fn(&result)), AfterNodeRollback(fn(&result)))

				return &result, arguments
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, options := test.arguments()
			client, err := New(options...)
			if err != nil {
				t.Fatal(err)
			}

			if client.Options.AfterNodeRollback == nil {
				t.Fatal("unexpected result expected option to be set")
			}

			_, err = client.Options.AfterNodeRollback(nil, nil, nil, nil, nil)
			if err != nil {
				t.Fatal(err)
			}

			if *result != test.expected {
				t.Fatalf("unexpected result %d, expected %d", *result, test.expected)
			}
		})
	}
}

func TestAfterNodeDoOption(t *testing.T) {
	fn := func(i *int) flow.AfterNodeHandler {
		return func(next flow.AfterNode) flow.AfterNode {
			return func(ctx context.Context, node *flow.Node, tracker *flow.Tracker, processes *flow.Processes, store refs.Store) (context.Context, error) {
				*i++
				return next(ctx, node, tracker, processes, store)
			}
		}
	}

	type test struct {
		expected  int
		arguments func() (*int, []api.Option)
	}

	tests := map[string]test{
		"single": {
			expected: 1,
			arguments: func() (*int, []api.Option) {
				result := 0
				arguments := NewCollection(AfterNodeDo(fn(&result)))

				return &result, arguments
			},
		},
		"multiple": {
			expected: 2,
			arguments: func() (*int, []api.Option) {
				result := 0
				arguments := NewCollection(AfterNodeDo(fn(&result)), AfterNodeDo(fn(&result)))

				return &result, arguments
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, options := test.arguments()
			client, err := New(options...)
			if err != nil {
				t.Fatal(err)
			}

			if client.Options.AfterNodeDo == nil {
				t.Fatal("unexpected result expected option to be set")
			}

			_, err = client.Options.AfterNodeDo(nil, nil, nil, nil, nil)
			if err != nil {
				t.Fatal(err)
			}

			if *result != test.expected {
				t.Fatalf("unexpected result %d, expected %d", *result, test.expected)
			}
		})
	}
}
