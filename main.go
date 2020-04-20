package maestro

import (
	"sync"

	"github.com/jexia/maestro/internal/constructor"
	"github.com/jexia/maestro/internal/logger"
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/transport"
)

// Client represents a maestro instance
type Client struct {
	Ctx       instance.Context
	Endpoints []*transport.Endpoint
	Flows     *specs.FlowsManifest
	Services  *specs.ServicesManifest
	Schema    *specs.SchemaManifest
	Listeners []transport.Listener
	Options   constructor.Options
}

// Serve opens all listeners inside the given maestro client
func (client *Client) Serve() (result error) {
	wg := sync.WaitGroup{}
	wg.Add(len(client.Listeners))

	for _, listener := range client.Listeners {
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

// Close gracefully closes the given client
func (client *Client) Close() {
	for _, listener := range client.Listeners {
		listener.Close()
	}

	for _, endpoint := range client.Endpoints {
		endpoint.Flow.Wait()
	}
}

// New constructs a new Maestro instance
func New(opts ...constructor.Option) (*Client, error) {
	ctx := instance.NewContext()
	options := NewOptions(ctx, opts...)

	mem, flows, endpoints, services, schema, err := constructor.Specs(ctx, options)
	if err != nil {
		return nil, err
	}

	managers, err := constructor.FlowManager(ctx, mem, services, endpoints, flows, options)
	if err != nil {
		return nil, err
	}

	client := &Client{
		Ctx:       ctx,
		Endpoints: managers,
		Flows:     flows,
		Services:  services,
		Schema:    schema,
		Listeners: options.Listeners,
		Options:   options,
	}

	return client, nil
}
