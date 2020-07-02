package maestro

import (
	"sync"

	"github.com/jexia/maestro/internal/core"
	"github.com/jexia/maestro/pkg/core/api"
	"github.com/jexia/maestro/pkg/core/instance"
	"github.com/jexia/maestro/pkg/core/logger"
	"github.com/jexia/maestro/pkg/core/trace"
	"github.com/jexia/maestro/pkg/functions"
	"github.com/jexia/maestro/pkg/transport"
)

// Client represents a maestro instance
type Client struct {
	Ctx          instance.Context
	transporters []*transport.Endpoint
	listeners    []transport.Listener
	collection   *api.Collection
	Options      api.Options
	mutex        sync.RWMutex
}

// Serve opens all listeners inside the given maestro client
func (client *Client) Serve() (result error) {
	if len(client.listeners) == 0 {
		return trace.New(trace.WithMessage("no listeners configured to serve"))
	}

	wg := sync.WaitGroup{}
	wg.Add(len(client.listeners))

	for _, listener := range client.listeners {
		client.Ctx.Logger(logger.Core).WithField("listener", listener.Name()).Info("serving listener")

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

// Handle updates the flows with the given specs collection.
// The given functions collection is used to execute functions on runtime.
func (client *Client) Handle(ctx instance.Context, options api.Options) error {
	client.mutex.Lock()
	defer client.mutex.Unlock()

	mem := functions.Collection{}
	collection, err := core.Specs(ctx, mem, options)
	if err != nil {
		return err
	}

	client.collection = collection
	managers, err := core.FlowManager(ctx, mem, collection.Services, collection.Endpoints, collection.Flows, options)
	if err != nil {
		return err
	}

	client.transporters = managers
	return nil
}

// Collection returns the currently defined specs collection
func (client *Client) Collection() *api.Collection {
	client.mutex.RLock()
	defer client.mutex.RUnlock()
	return client.collection
}

// Close gracefully closes the given client
func (client *Client) Close() {
	for _, listener := range client.listeners {
		listener.Close()
	}

	for _, transporter := range client.transporters {
		if transporter.Flow == nil {
			continue
		}

		transporter.Flow.Wait()
	}
}

// New constructs a new Maestro instance
func New(opts ...api.Option) (*Client, error) {
	ctx := instance.NewContext()
	options, err := NewOptions(ctx, opts...)
	if err != nil {
		return nil, err
	}

	client := &Client{
		Ctx:       ctx,
		listeners: options.Listeners,
		Options:   options,
	}

	err = client.Handle(ctx, options)
	if err != nil {
		return nil, err
	}

	return client, nil
}
