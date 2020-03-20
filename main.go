package maestro

import (
	"sync"

	"github.com/jexia/maestro/constructor"
	"github.com/jexia/maestro/protocol"
	"github.com/jexia/maestro/specs"
	log "github.com/sirupsen/logrus"
)

// Client represents a maestro instance
type Client struct {
	Endpoints []*protocol.Endpoint
	Manifest  *specs.Manifest
	Listeners []protocol.Listener
	Options   constructor.Options
}

// Serve opens all listeners inside the given maestro client
func (client *Client) Serve() (result error) {
	wg := sync.WaitGroup{}
	wg.Add(len(client.Listeners))

	for _, listener := range client.Listeners {
		log.WithField("listener", listener.Name()).Info("serving listener")

		go func(listener protocol.Listener) {
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
	options := constructor.NewOptions(opts...)

	manifest, err := constructor.Specs(options)
	if err != nil {
		return nil, err
	}

	endpoints, err := constructor.FlowManager(manifest, options)
	if err != nil {
		return nil, err
	}

	client := &Client{
		Endpoints: endpoints,
		Manifest:  manifest,
		Listeners: options.Listeners,
		Options:   options,
	}

	return client, nil
}
