package daemon

import (
	"errors"
	"sync"

	"github.com/jexia/semaphore"
	"github.com/jexia/semaphore/cmd/semaphore/daemon/providers"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/endpoints"
	"github.com/jexia/semaphore/pkg/broker/listeners"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/broker/trace"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/transport"
	"go.uber.org/zap"
)

// NewClient constructs a new Semaphore instance
func NewClient(ctx *broker.Context, core semaphore.Options, provider providers.Options) (*Client, error) {

	if ctx == nil {
		return nil, errors.New("nil context")
	}

	client := &Client{
		ctx:       ctx,
		core:      core,
		providers: provider,
		stack:     functions.Collection{},
	}

	err := client.Apply(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Client represents a semaphore instance
type Client struct {
	core      semaphore.Options
	providers providers.Options
	ctx       *broker.Context
	mutex     sync.Mutex
	stack     functions.Collection
}

// Apply updates the listeners with the given specs collection.
// Transporters are created from the available endpoints and flows.
// The created transporters are passed to the listeners to be hot-swapped.
//
// This method does not perform any checks ensuring that the given
// specification is valid.
func (client *Client) Apply(ctx *broker.Context) error {
	client.mutex.Lock()
	defer client.mutex.Unlock()

	collection, err := providers.Resolve(ctx, client.stack, client.providers)
	if err != nil {
		return err
	}

	transporters, err := endpoints.Transporters(ctx, collection.EndpointList, collection.FlowListInterface,
		endpoints.WithCore(client.core),
		endpoints.WithServices(collection.ServiceList),
		endpoints.WithFunctions(client.stack),
	)

	if err != nil {
		return err
	}

	err = listeners.Apply(ctx, client.providers.Codec, client.providers.Listeners, transporters)
	if err != nil {
		return err
	}

	return nil
}

// Serve opens all listeners inside the given semaphore client
func (client *Client) Serve() (result error) {
	if len(client.providers.Listeners) == 0 {
		return trace.New(trace.WithMessage("no listeners configured to serve"))
	}

	wg := sync.WaitGroup{}
	wg.Add(len(client.providers.Listeners))

	for _, listener := range client.providers.Listeners {
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
	for _, listener := range client.providers.Listeners {
		listener.Close()
	}
}
