package maestro

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/jexia/maestro/definitions/hcl"
	"github.com/jexia/maestro/protocol"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/strict"
	"github.com/jexia/maestro/specs/trace"
	"github.com/jexia/maestro/utils"
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

	if options.Path == "" {
		return nil, trace.New(trace.WithMessage("undefined path in options"))
	}

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
	files, err := utils.ReadDir(options.Path, options.Recursive, hcl.Ext)
	if err != nil {
		return nil, err
	}

	manifest := &specs.Manifest{}

	for _, file := range files {
		reader, err := os.Open(filepath.Join(file.Path, file.Name()))
		if err != nil {
			return nil, err
		}

		definition, err := hcl.UnmarshalHCL(file.Name(), reader)
		if err != nil {
			return nil, err
		}

		result, err := hcl.ParseSpecs(definition, options.Functions)
		if err != nil {
			return nil, err
		}

		collection, err := hcl.ParseSchema(definition, options.Schema)
		if err != nil {
			return nil, err
		}

		options.Schema.Add(collection)
		manifest.MergeLeft(result)

		err = specs.CheckManifestDuplicates(file.Name(), manifest)
		if err != nil {
			return nil, err
		}
	}

	err = specs.ResolveManifestDependencies(manifest)
	if err != nil {
		return nil, err
	}

	err = strict.DefineManifest(options.Schema, manifest)
	if err != nil {
		return nil, err
	}

	return manifest, nil
}
