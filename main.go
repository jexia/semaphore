package semaphore

import (
	"sync"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/config"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/core"
	"github.com/jexia/semaphore/pkg/core/trace"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/transport"
	"go.uber.org/zap"
)

// Client represents a semaphore instance
type Client struct {
	Ctx          *broker.Context
	transporters []*transport.Endpoint
	listeners    []transport.Listener
	flows        specs.FlowListInterface
	endpoints    specs.EndpointList
	services     specs.ServiceList
	schemas      specs.Schemas
	Options      config.Options
	mutex        sync.RWMutex
}

// Serve opens all listeners inside the given semaphore client
func (client *Client) Serve() (result error) {
	if len(client.listeners) == 0 {
		return trace.New(trace.WithMessage("no listeners configured to serve"))
	}

	wg := sync.WaitGroup{}
	wg.Add(len(client.listeners))

	for _, listener := range client.listeners {
		logger.Info(client.Ctx, "serving listener", zap.String("listener", listener.Name()))

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
func (client *Client) Handle(ctx *broker.Context, options config.Options) error {
	client.mutex.Lock()
	defer client.mutex.Unlock()

	mem := functions.Collection{}
	flows, endpoints, services, schemas, err := options.Constructor(ctx, mem, options)
	if err != nil {
		return err
	}

	managers, err := core.Apply(ctx, mem, services, endpoints, flows, options)
	if err != nil {
		return err
	}

	client.flows = flows
	client.endpoints = endpoints
	client.services = services
	client.schemas = schemas
	client.transporters = managers

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
func New(opts ...config.Option) (*Client, error) {
	ctx := logger.WithLogger(broker.NewContext())
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
