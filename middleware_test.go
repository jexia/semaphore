package semaphore

import (
	"context"
	"testing"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/flow"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

func TestWithMiddleware(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	middleware := MiddlewareFunc(func(*broker.Context) ([]Option, error) {
		return nil, nil
	})

	client, err := NewOptions(ctx, WithMiddleware(middleware), WithMiddleware(middleware))
	if err != nil {
		t.Fatal(err)
	}

	if client.Middleware == nil {
		t.Fatal("middleware not set")
	}

	if len(client.Middleware) != 2 {
		t.Fatalf("unexpected middleware %d, expected 2", len(client.Middleware))
	}
}

func TestBeforeManagerDoOption(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())

	fn := func(i *int) flow.BeforeManagerHandler {
		return func(next flow.BeforeManager) flow.BeforeManager {
			return func(ctx context.Context, manager *flow.Manager, store references.Store) (context.Context, error) {
				*i++
				return next(ctx, manager, store)
			}
		}
	}

	type test struct {
		expected  int
		arguments func() (*int, []Option)
	}

	tests := map[string]test{
		"single": {
			expected: 1,
			arguments: func() (*int, []Option) {
				result := 0
				arguments := NewCollection(BeforeManagerDo(fn(&result)))

				return &result, arguments
			},
		},
		"multiple": {
			expected: 2,
			arguments: func() (*int, []Option) {
				result := 0
				arguments := NewCollection(BeforeManagerDo(fn(&result)), BeforeManagerDo(fn(&result)))

				return &result, arguments
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, options := test.arguments()
			client, err := NewOptions(ctx, options...)
			if err != nil {
				t.Fatal(err)
			}

			if client.BeforeManagerDo == nil {
				t.Fatal("unexpected result expected option to be set")
			}

			_, err = client.BeforeManagerDo(nil, nil, nil)
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
	ctx := logger.WithLogger(broker.NewBackground())

	fn := func(i *int) flow.BeforeManagerHandler {
		return func(next flow.BeforeManager) flow.BeforeManager {
			return func(ctx context.Context, manager *flow.Manager, store references.Store) (context.Context, error) {
				*i++
				return next(ctx, manager, store)
			}
		}
	}

	type test struct {
		expected  int
		arguments func() (*int, []Option)
	}

	tests := map[string]test{
		"single": {
			expected: 1,
			arguments: func() (*int, []Option) {
				result := 0
				arguments := NewCollection(BeforeManagerRollback(fn(&result)))

				return &result, arguments
			},
		},
		"multiple": {
			expected: 2,
			arguments: func() (*int, []Option) {
				result := 0
				arguments := NewCollection(BeforeManagerRollback(fn(&result)), BeforeManagerRollback(fn(&result)))

				return &result, arguments
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, options := test.arguments()
			client, err := NewOptions(ctx, options...)
			if err != nil {
				t.Fatal(err)
			}

			if client.BeforeManagerRollback == nil {
				t.Fatal("unexpected result expected option to be set")
			}

			_, err = client.BeforeManagerRollback(nil, nil, nil)
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
	ctx := logger.WithLogger(broker.NewBackground())

	fn := func(i *int) flow.AfterManagerHandler {
		return func(next flow.AfterManager) flow.AfterManager {
			return func(ctx context.Context, manager *flow.Manager, store references.Store) (context.Context, error) {
				*i++
				return next(ctx, manager, store)
			}
		}
	}

	type test struct {
		expected  int
		arguments func() (*int, []Option)
	}

	tests := map[string]test{
		"single": {
			expected: 1,
			arguments: func() (*int, []Option) {
				result := 0
				arguments := NewCollection(AfterManagerRollback(fn(&result)))

				return &result, arguments
			},
		},
		"multiple": {
			expected: 2,
			arguments: func() (*int, []Option) {
				result := 0
				arguments := NewCollection(AfterManagerRollback(fn(&result)), AfterManagerRollback(fn(&result)))

				return &result, arguments
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, options := test.arguments()
			client, err := NewOptions(ctx, options...)
			if err != nil {
				t.Fatal(err)
			}

			if client.AfterManagerRollback == nil {
				t.Fatal("unexpected result expected option to be set")
			}

			_, err = client.AfterManagerRollback(nil, nil, nil)
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
	ctx := logger.WithLogger(broker.NewBackground())

	fn := func(i *int) flow.AfterManagerHandler {
		return func(next flow.AfterManager) flow.AfterManager {
			return func(ctx context.Context, manager *flow.Manager, store references.Store) (context.Context, error) {
				*i++
				return next(ctx, manager, store)
			}
		}
	}

	type test struct {
		expected  int
		arguments func() (*int, []Option)
	}

	tests := map[string]test{
		"single": {
			expected: 1,
			arguments: func() (*int, []Option) {
				result := 0
				arguments := NewCollection(AfterManagerDo(fn(&result)))

				return &result, arguments
			},
		},
		"multiple": {
			expected: 2,
			arguments: func() (*int, []Option) {
				result := 0
				arguments := NewCollection(AfterManagerDo(fn(&result)), AfterManagerDo(fn(&result)))

				return &result, arguments
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, options := test.arguments()
			client, err := NewOptions(ctx, options...)
			if err != nil {
				t.Fatal(err)
			}

			if client.AfterManagerDo == nil {
				t.Fatal("unexpected result expected option to be set")
			}

			_, err = client.AfterManagerDo(nil, nil, nil)
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
	ctx := logger.WithLogger(broker.NewBackground())

	fn := func(i *int) flow.BeforeNodeHandler {
		return func(next flow.BeforeNode) flow.BeforeNode {
			return func(ctx context.Context, node *flow.Node, tracker flow.Tracker, processes *flow.Processes, store references.Store) (context.Context, error) {
				*i++
				return next(ctx, node, tracker, processes, store)
			}
		}
	}

	type test struct {
		expected  int
		arguments func() (*int, []Option)
	}

	tests := map[string]test{
		"single": {
			expected: 1,
			arguments: func() (*int, []Option) {
				result := 0
				arguments := NewCollection(BeforeNodeDo(fn(&result)))

				return &result, arguments
			},
		},
		"multiple": {
			expected: 2,
			arguments: func() (*int, []Option) {
				result := 0
				arguments := NewCollection(BeforeNodeDo(fn(&result)), BeforeNodeDo(fn(&result)))

				return &result, arguments
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, options := test.arguments()
			client, err := NewOptions(ctx, options...)
			if err != nil {
				t.Fatal(err)
			}

			if client.BeforeNodeDo == nil {
				t.Fatal("unexpected result expected option to be set")
			}

			_, err = client.BeforeNodeDo(nil, nil, nil, nil, nil)
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
	ctx := logger.WithLogger(broker.NewBackground())

	fn := func(i *int) flow.BeforeNodeHandler {
		return func(next flow.BeforeNode) flow.BeforeNode {
			return func(ctx context.Context, node *flow.Node, tracker flow.Tracker, processes *flow.Processes, store references.Store) (context.Context, error) {
				*i++
				return next(ctx, node, tracker, processes, store)
			}
		}
	}

	type test struct {
		expected  int
		arguments func() (*int, []Option)
	}

	tests := map[string]test{
		"single": {
			expected: 1,
			arguments: func() (*int, []Option) {
				result := 0
				arguments := NewCollection(BeforeNodeRollback(fn(&result)))

				return &result, arguments
			},
		},
		"multiple": {
			expected: 2,
			arguments: func() (*int, []Option) {
				result := 0
				arguments := NewCollection(BeforeNodeRollback(fn(&result)), BeforeNodeRollback(fn(&result)))

				return &result, arguments
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, options := test.arguments()
			client, err := NewOptions(ctx, options...)
			if err != nil {
				t.Fatal(err)
			}

			if client.BeforeNodeRollback == nil {
				t.Fatal("unexpected result expected option to be set")
			}

			_, err = client.BeforeNodeRollback(nil, nil, nil, nil, nil)
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
	ctx := logger.WithLogger(broker.NewBackground())

	fn := func(i *int) flow.AfterNodeHandler {
		return func(next flow.AfterNode) flow.AfterNode {
			return func(ctx context.Context, node *flow.Node, tracker flow.Tracker, processes *flow.Processes, store references.Store) (context.Context, error) {
				*i++
				return next(ctx, node, tracker, processes, store)
			}
		}
	}

	type test struct {
		expected  int
		arguments func() (*int, []Option)
	}

	tests := map[string]test{
		"single": {
			expected: 1,
			arguments: func() (*int, []Option) {
				result := 0
				arguments := NewCollection(AfterNodeRollback(fn(&result)))

				return &result, arguments
			},
		},
		"multiple": {
			expected: 2,
			arguments: func() (*int, []Option) {
				result := 0
				arguments := NewCollection(AfterNodeRollback(fn(&result)), AfterNodeRollback(fn(&result)))

				return &result, arguments
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, options := test.arguments()
			client, err := NewOptions(ctx, options...)
			if err != nil {
				t.Fatal(err)
			}

			if client.AfterNodeRollback == nil {
				t.Fatal("unexpected result expected option to be set")
			}

			_, err = client.AfterNodeRollback(nil, nil, nil, nil, nil)
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
	ctx := logger.WithLogger(broker.NewBackground())

	fn := func(i *int) flow.AfterNodeHandler {
		return func(next flow.AfterNode) flow.AfterNode {
			return func(ctx context.Context, node *flow.Node, tracker flow.Tracker, processes *flow.Processes, store references.Store) (context.Context, error) {
				*i++
				return next(ctx, node, tracker, processes, store)
			}
		}
	}

	type test struct {
		expected  int
		arguments func() (*int, []Option)
	}

	tests := map[string]test{
		"single": {
			expected: 1,
			arguments: func() (*int, []Option) {
				result := 0
				arguments := NewCollection(AfterNodeDo(fn(&result)))

				return &result, arguments
			},
		},
		"multiple": {
			expected: 2,
			arguments: func() (*int, []Option) {
				result := 0
				arguments := NewCollection(AfterNodeDo(fn(&result)), AfterNodeDo(fn(&result)))

				return &result, arguments
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, options := test.arguments()
			client, err := NewOptions(ctx, options...)
			if err != nil {
				t.Fatal(err)
			}

			if client.AfterNodeDo == nil {
				t.Fatal("unexpected result expected option to be set")
			}

			_, err = client.AfterNodeDo(nil, nil, nil, nil, nil)
			if err != nil {
				t.Fatal(err)
			}

			if *result != test.expected {
				t.Fatalf("unexpected result %d, expected %d", *result, test.expected)
			}
		})
	}
}

func TestBeforeConstructor(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())

	fn := func(i *int) BeforeConstructorHandler {
		return func(next BeforeConstructor) BeforeConstructor {
			return func(ctx *broker.Context, fns functions.Collection, options Options) error {
				*i++
				return next(ctx, fns, options)
			}
		}
	}

	type test struct {
		expected  int
		arguments func() (*int, []Option)
	}

	tests := map[string]test{
		"single": {
			expected: 1,
			arguments: func() (*int, []Option) {
				result := 0
				arguments := NewCollection(WithBeforeConstructor(fn(&result)))

				return &result, arguments
			},
		},
		"multiple": {
			expected: 2,
			arguments: func() (*int, []Option) {
				result := 0
				arguments := NewCollection(WithBeforeConstructor(fn(&result)), WithBeforeConstructor(fn(&result)))

				return &result, arguments
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, input := test.arguments()
			options, err := NewOptions(ctx, input...)
			if err != nil {
				t.Fatal(err)
			}

			if options.BeforeConstructor == nil {
				t.Fatal("unexpected result expected option to be set")
			}

			err = options.BeforeConstructor(nil, nil, Options{})
			if err != nil {
				t.Fatal(err)
			}

			if *result != test.expected {
				t.Fatalf("unexpected result %d, expected %d", *result, test.expected)
			}
		})
	}
}

func TestAfterFlowConstructor(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())

	fn := func(i *int) AfterFlowConstructionHandler {
		return func(next AfterFlowConstruction) AfterFlowConstruction {
			return func(ctx *broker.Context, flow specs.FlowInterface, manager *flow.Manager) error {
				*i++
				return next(ctx, flow, manager)
			}
		}
	}

	type test struct {
		expected  int
		arguments func() (*int, []Option)
	}

	tests := map[string]test{
		"single": {
			expected: 1,
			arguments: func() (*int, []Option) {
				result := 0
				arguments := NewCollection(AfterFlowConstructor(fn(&result)))

				return &result, arguments
			},
		},
		"multiple": {
			expected: 2,
			arguments: func() (*int, []Option) {
				result := 0
				arguments := NewCollection(AfterFlowConstructor(fn(&result)), AfterFlowConstructor(fn(&result)))

				return &result, arguments
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, input := test.arguments()
			options, err := NewOptions(ctx, input...)
			if err != nil {
				t.Fatal(err)
			}

			if options.AfterFlowConstruction == nil {
				t.Fatal("unexpected result expected option to be set")
			}

			err = options.AfterFlowConstruction(nil, nil, nil)
			if err != nil {
				t.Fatal(err)
			}

			if *result != test.expected {
				t.Fatalf("unexpected result %d, expected %d", *result, test.expected)
			}
		})
	}
}
