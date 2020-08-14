package semaphore

import (
	"errors"
	"sync"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/config"
	"github.com/jexia/semaphore/pkg/broker/endpoints"
	"github.com/jexia/semaphore/pkg/broker/listeners"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/broker/providers"
	"github.com/jexia/semaphore/pkg/broker/trace"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/transport"
	"go.uber.org/zap"
)

// Client represents a semaphore instance
type Client struct {
	config.Options
	ctx          *broker.Context
	transporters transport.EndpointList
	listeners    transport.ListenerList
	flows        specs.FlowListInterface
	endpoints    specs.EndpointList
	services     specs.ServiceList
	schemas      specs.Schemas
	mutex        sync.RWMutex
	stack        functions.Collection
}

// Serve opens all listeners inside the given semaphore client
func (client *Client) Serve() (result error) {
	if len(client.listeners) == 0 {
		return trace.New(trace.WithMessage("no listeners configured to serve"))
	}

	wg := sync.WaitGroup{}
	wg.Add(len(client.listeners))

	for _, listener := range client.listeners {
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

// Resolve resolves the configured providers and constructs a valid Semaphore
// specification. Any error thrown during the compilation of the specification
// is returned. Type checks are performed after the specifications have been
// resolved ensuring type safety.
func (client *Client) Resolve(ctx *broker.Context) (providers.Collection, error) {
	return providers.Resolve(ctx, client.stack, client.Options)
}

// Apply updates the listeners with the given specs collection.
// Transporters are created from the available endpoints and flows.
// The created transporters are passed to the listeners to be hot-swapped.
//
// This method does not perform any checks checking ensuring strict types
// or whether the specification is valid.
func (client *Client) Apply(ctx *broker.Context, collection providers.Collection) error {
	client.mutex.Lock()
	defer client.mutex.Unlock()

	transporters, err := endpoints.Transporters(ctx, collection.EndpointList, collection.FlowListInterface,
		endpoints.WithServices(collection.ServiceList),
		endpoints.WithOptions(client.Options),
		endpoints.WithFunctions(client.stack),
	)

	if err != nil {
		return err
	}

	err = listeners.Apply(ctx, client.Codec, client.listeners, transporters)
	if err != nil {
		return err
	}

	client.flows = collection.FlowListInterface
	client.endpoints = collection.EndpointList
	client.services = collection.ServiceList
	client.schemas = collection.Schemas
	client.transporters = transporters

	return nil
}

// GetFlows returns the currently applied flows
func (client *Client) GetFlows() specs.FlowListInterface {
	return client.flows
}

// GetServices returns the currently applied services
func (client *Client) GetServices() specs.ServiceList {
	return client.services
}

// GetEndpoints returns the currently applied endpoints
func (client *Client) GetEndpoints() specs.EndpointList {
	return client.endpoints
}

// GetSchemas returns the currently applied schemas
func (client *Client) GetSchemas() specs.Schemas {
	return client.schemas
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

// New constructs a new Semaphore instance
func New(ctx *broker.Context, opts ...config.Option) (*Client, error) {
	if ctx == nil {
		return nil, errors.New("nil context")
	}

	options, err := NewOptions(ctx, opts...)
	if err != nil {
		return nil, err
	}

	client := &Client{
		ctx:       ctx,
		listeners: options.Listeners,
		Options:   options,
		stack:     functions.Collection{},
	}

	specs, err := client.Resolve(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Apply(ctx, specs)
	if err != nil {
		return nil, err
	}

	return client, nil
}
