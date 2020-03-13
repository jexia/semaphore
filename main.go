package maestro

import (
	"sync"

	"github.com/jexia/maestro/protocol"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/strict"
	log "github.com/sirupsen/logrus"
)

// Client represents a maestro instance
type Client struct {
	Endpoints []*protocol.Endpoint
	Manifest  *specs.Manifest
	Listeners []protocol.Listener
	Options   Options
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
func New(opts ...Option) (*Client, error) {
	options := NewOptions(opts...)

	manifest, err := ConstructSpecs(options)
	if err != nil {
		return nil, err
	}

	endpoints, err := ConstructFlowManager(manifest, options)
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

// ConstructSpecs construct a specs manifest from the given options
func ConstructSpecs(options Options) (*specs.Manifest, error) {
	result := &specs.Manifest{}

	for _, resolver := range options.Definitions {
		if resolver == nil {
			continue
		}

		manifest, err := resolver(options.Functions)
		if err != nil {
			return nil, err
		}

		result.Merge(manifest)
	}

	for _, resolver := range options.Schemas {
		if resolver == nil {
			continue
		}

		err := resolver(options.Schema)
		if err != nil {
			return nil, err
		}
	}

	err := specs.CheckManifestDuplicates(result)
	if err != nil {
		return nil, err
	}

	err = specs.ResolveManifestDependencies(result)
	if err != nil {
		return nil, err
	}

	err = strict.DefineManifest(options.Schema, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
