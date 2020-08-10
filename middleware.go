package semaphore

import (
	"context"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/config"
	"github.com/jexia/semaphore/pkg/flow"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

// WithMiddleware initialises the given middleware and defines all options
func WithMiddleware(middleware config.Middleware) config.Option {
	return func(options *config.Options) {
		options.Middleware = append(options.Middleware, middleware)
	}
}

// BeforeConstructor the passed function gets called before new specifications are constructed
func BeforeConstructor(wrapper config.BeforeConstructorHandler) config.Option {
	return func(options *config.Options) {
		if options.BeforeConstructor == nil {
			options.BeforeConstructor = wrapper(func(*broker.Context, functions.Collection, config.Options) error { return nil })
			return
		}

		options.BeforeConstructor = wrapper(options.BeforeConstructor)
	}
}

// AfterConstructor the passed function gets called once all options have been applied
func AfterConstructor(wrapper config.AfterConstructorHandler) config.Option {
	return func(options *config.Options) {
		if options.AfterConstructor == nil {
			options.AfterConstructor = wrapper(func(*broker.Context, specs.FlowListInterface, specs.EndpointList, specs.ServiceList, specs.Schemas) error {
				return nil
			})
			return
		}

		options.AfterConstructor = wrapper(options.AfterConstructor)
	}
}

// BeforeManagerDo the passed function gets called before a request gets handled by a flow manager
func BeforeManagerDo(wrapper flow.BeforeManagerHandler) config.Option {
	return func(options *config.Options) {
		if options.BeforeManagerDo == nil {
			options.BeforeManagerDo = wrapper(func(ctx context.Context, manager *flow.Manager, store references.Store) (context.Context, error) {
				return ctx, nil
			})

			return
		}

		options.BeforeManagerDo = wrapper(options.BeforeManagerDo)
	}
}

// BeforeManagerRollback the passed function gets called before a rollback request gets handled by a flow manager
func BeforeManagerRollback(wrapper flow.BeforeManagerHandler) config.Option {
	return func(options *config.Options) {
		if options.BeforeManagerRollback == nil {
			options.BeforeManagerRollback = wrapper(func(ctx context.Context, manager *flow.Manager, store references.Store) (context.Context, error) {
				return ctx, nil
			})

			return
		}

		options.BeforeManagerRollback = wrapper(options.BeforeManagerRollback)
	}
}

// AfterManagerDo the passed function gets after a flow has been handled by the flow manager
func AfterManagerDo(wrapper flow.AfterManagerHandler) config.Option {
	return func(options *config.Options) {
		if options.AfterManagerDo == nil {
			options.AfterManagerDo = wrapper(func(ctx context.Context, manager *flow.Manager, store references.Store) (context.Context, error) {
				return ctx, nil
			})

			return
		}

		options.AfterManagerDo = wrapper(options.AfterManagerDo)
	}
}

// AfterManagerRollback the passed function gets after a flow rollback has been handled by the flow manager
func AfterManagerRollback(wrapper flow.AfterManagerHandler) config.Option {
	return func(options *config.Options) {
		if options.AfterManagerRollback == nil {
			options.AfterManagerRollback = wrapper(func(ctx context.Context, manager *flow.Manager, store references.Store) (context.Context, error) {
				return ctx, nil
			})

			return
		}

		options.AfterManagerRollback = wrapper(options.AfterManagerRollback)
	}
}

// BeforeNodeDo the passed function gets called before a node is executed
func BeforeNodeDo(wrapper flow.BeforeNodeHandler) config.Option {
	return func(options *config.Options) {
		if options.BeforeNodeDo == nil {
			options.BeforeNodeDo = wrapper(func(ctx context.Context, node *flow.Node, tracker *flow.Tracker, processes *flow.Processes, store references.Store) (context.Context, error) {
				return ctx, nil
			})

			return
		}

		options.BeforeNodeDo = wrapper(options.BeforeNodeDo)
	}
}

// BeforeNodeRollback the passed function gets called before a node rollback is executed
func BeforeNodeRollback(wrapper flow.BeforeNodeHandler) config.Option {
	return func(options *config.Options) {
		if options.BeforeNodeRollback == nil {
			options.BeforeNodeRollback = wrapper(func(ctx context.Context, node *flow.Node, tracker *flow.Tracker, processes *flow.Processes, store references.Store) (context.Context, error) {
				return ctx, nil
			})

			return
		}

		options.BeforeNodeRollback = wrapper(options.BeforeNodeRollback)
	}
}

// AfterNodeDo the passed function gets called after a node is executed
func AfterNodeDo(wrapper flow.AfterNodeHandler) config.Option {
	return func(options *config.Options) {
		if options.AfterNodeDo == nil {
			options.AfterNodeDo = wrapper(func(ctx context.Context, node *flow.Node, tracker *flow.Tracker, processes *flow.Processes, store references.Store) (context.Context, error) {
				return ctx, nil
			})
			return
		}

		options.AfterNodeDo = wrapper(options.AfterNodeDo)
	}
}

// AfterNodeRollback the passed function gets called after a node rollback is executed
func AfterNodeRollback(wrapper flow.AfterNodeHandler) config.Option {
	return func(options *config.Options) {
		if options.AfterNodeRollback == nil {
			options.AfterNodeRollback = wrapper(func(ctx context.Context, node *flow.Node, tracker *flow.Tracker, processes *flow.Processes, store references.Store) (context.Context, error) {
				return ctx, nil
			})
			return
		}

		options.AfterNodeRollback = wrapper(options.AfterNodeRollback)
	}
}
