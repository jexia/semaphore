package maestro

import (
	"sync"

	"github.com/jexia/maestro/internal/constructor"
	"github.com/jexia/maestro/internal/functions"
	"github.com/jexia/maestro/pkg/core/api"
	"github.com/jexia/maestro/pkg/core/instance"
	"github.com/jexia/maestro/pkg/core/logger"
	"github.com/jexia/maestro/pkg/core/trace"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/transport"
)

// Client represents a maestro instance
type Client struct {
	Ctx          instance.Context
	Transporters []*transport.Endpoint
	Flows        *specs.FlowsManifest
	Services     *specs.ServicesManifest
	Schema       *specs.SchemaManifest
	Endpoints    *specs.EndpointsManifest
	Listeners    []transport.Listener
	Options      api.Options
}

// Serve opens all listeners inside the given maestro client
func (client *Client) Serve() (result error) {
	if len(client.Listeners) == 0 {
		return trace.New(trace.WithMessage("no listeners configured to serve"))
	}

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

	for _, transporter := range client.Transporters {
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

	mem := functions.Collection{}
	collection, err := constructor.Specs(ctx, mem, options)
	if err != nil {
		return nil, err
	}

	managers, err := constructor.FlowManager(ctx, mem, collection.Services, collection.Endpoints, collection.Flows, options)
	if err != nil {
		return nil, err
	}

	client := &Client{
		Ctx:          ctx,
		Transporters: managers,
		Flows:        collection.Flows,
		Services:     collection.Services,
		Schema:       collection.Schema,
		Endpoints:    collection.Endpoints,
		Listeners:    options.Listeners,
		Options:      options,
	}

	return client, nil
}
