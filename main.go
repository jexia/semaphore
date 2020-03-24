package maestro

import (
	"context"
	"sync"

	"github.com/jexia/maestro/constructor"
	"github.com/jexia/maestro/logger"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/transport"
)

// Client represents a maestro instance
type Client struct {
	ctx       context.Context
	Endpoints []*transport.Endpoint
	Manifest  *specs.Manifest
	Listeners []transport.Listener
	Options   constructor.Options
}

// Serve opens all listeners inside the given maestro client
func (client *Client) Serve() (result error) {
	wg := sync.WaitGroup{}
	wg.Add(len(client.Listeners))

	for _, listener := range client.Listeners {
		logger.FromCtx(client.ctx, logger.Core).WithField("listener", listener.Name()).Info("serving listener")

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
	ctx := context.Background()
	ctx = logger.WithValue(ctx)

	options := constructor.NewOptions(ctx, opts...)

	manifest, err := constructor.Specs(ctx, options)
	if err != nil {
		return nil, err
	}

	endpoints, err := constructor.FlowManager(ctx, manifest, options)
	if err != nil {
		return nil, err
	}

	client := &Client{
		ctx:       ctx,
		Endpoints: endpoints,
		Manifest:  manifest,
		Listeners: options.Listeners,
		Options:   options,
	}

	return client, nil
}
