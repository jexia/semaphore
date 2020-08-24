package daemon

import (
	"errors"
	"sync"

	"github.com/jexia/semaphore"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/endpoints"
	"github.com/jexia/semaphore/pkg/broker/listeners"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/broker/providers"
	"github.com/jexia/semaphore/pkg/broker/trace"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/transport"
	"go.uber.org/zap"
)

// New constructs a new Semaphore instance
func New(ctx *broker.Context, opts ...semaphore.Option) (*Client, error) {
	if ctx == nil {
		return nil, errors.New("nil context")
	}

	// TODO: refactor client

	options, err := semaphore.NewOptions(ctx, opts...)
	if err != nil {
		return nil, err
	}

	mem := functions.Collection{}
	collection, err := providers.Resolve(ctx, mem, options)
	if err != nil {
		return nil, err
	}

	client := &Client{
		ctx:     ctx,
		Options: options,
		Stack:   mem,
	}

	err = client.Apply(ctx, collection)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Client represents a semaphore instance
type Client struct {
	semaphore.Options
	ctx   *broker.Context
	mutex sync.Mutex
	Stack functions.Collection
}

// Apply updates the listeners with the given specs collection.
// Transporters are created from the available endpoints and flows.
// The created transporters are passed to the listeners to be hot-swapped.
//
// This method does not perform any checks ensuring that the given
// specification is valid.
func (client *Client) Apply(ctx *broker.Context, collection providers.Collection) error {
	client.mutex.Lock()
	defer client.mutex.Unlock()

	transporters, err := endpoints.Transporters(ctx, collection.EndpointList, collection.FlowListInterface,
		endpoints.WithServices(collection.ServiceList),
		endpoints.WithOptions(client.Options),
		endpoints.WithFunctions(client.Stack),
	)

	if err != nil {
		return err
	}

	err = listeners.Apply(ctx, client.Codec, client.Options.Listeners, transporters)
	if err != nil {
		return err
	}

	return nil
}

// Serve opens all listeners inside the given semaphore client
func (client *Client) Serve() (result error) {
	if len(client.Options.Listeners) == 0 {
		return trace.New(trace.WithMessage("no listeners configured to serve"))
	}

	wg := sync.WaitGroup{}
	wg.Add(len(client.Options.Listeners))

	for _, listener := range client.Options.Listeners {
		logger.Info(client.ctx, "serving listener", zap.String("listener", listener.Name()))

		go func(listener transport.Listener) {
			defer wg.Done()
			err := listener.Serve()
			if err != nil {
				result = err
			}
		}(listener)
	}

	wg.Wait()
	return result
}

// Close gracefully closes the given client
func (client *Client) Close() {
	for _, listener := range client.Options.Listeners {
		listener.Close()
	}
}
