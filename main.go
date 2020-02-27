package maestro

import (
	"context"
	"os"
	"path/filepath"

	"github.com/jexia/maestro/flow"
	"github.com/jexia/maestro/refs"

	"github.com/jexia/maestro/codec"
	"github.com/jexia/maestro/protocol"
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/intermediate"
	"github.com/jexia/maestro/specs/strict"
	"github.com/jexia/maestro/specs/trace"
	"github.com/jexia/maestro/utils"
)

// Client represents a maestro instance
type Client struct {
	Manifest  *specs.Manifest
	Listeners []protocol.Listener
	Options   Options
}

// Serve opens the client listeners
func (client *Client) Serve() {

}

// Option represents a constructor func which sets a given option
type Option func(*Options)

// Options represents all the available options
type Options struct {
	Path      string
	Recursive bool
	Codec     map[string]codec.Constructor
	Callers   []protocol.Caller
	Listeners []protocol.Listener
	Schema    schema.Collection
	Functions specs.CustomDefinedFunctions
}

// NewOptions constructs a options object from the given option constructors
func NewOptions(options ...Option) Options {
	result := Options{}
	for _, option := range options {
		option(&result)
	}
	return result
}

// WithPath defines the definitions path to be used
func WithPath(path string, recursive bool) Option {
	return func(options *Options) {
		options.Path = path
		options.Recursive = recursive
	}
}

// WithCodec appends the given codec to the collection of available codecs
func WithCodec(constructor codec.Constructor) Option {
	return func(options *Options) {
		options.Codec[constructor.Name()] = constructor
	}
}

// WithSchemaCollection defines the schema collection to be used
func WithSchemaCollection(collection schema.Collection) Option {
	return func(options *Options) {
		options.Schema = collection
	}
}

// WithFunctions defines the custom defined functions to be used
func WithFunctions(functions specs.CustomDefinedFunctions) Option {
	return func(options *Options) {
		options.Functions = functions
	}
}

// New constructs a new Maestro instance
func New(opts ...Option) (*Client, error) {
	options := NewOptions(opts...)

	if options.Path == "" {
		return nil, trace.New(trace.WithMessage("undefined path in options"))
	}

	if options.Schema == nil {
		return nil, trace.New(trace.WithMessage("undefined schema in options"))
	}

	manifest, err := ConstructSpecs(options)
	if err != nil {
		return nil, err
	}

	err = ConstructFlowManager(manifest, options)
	if err != nil {
		return nil, err
	}

	client := &Client{
		Manifest:  manifest,
		Listeners: options.Listeners,
		Options:   options,
	}

	return client, nil
}

// ConstructSpecs construct a specs manifest from the given options
func ConstructSpecs(options Options) (*specs.Manifest, error) {
	files, err := utils.ReadDir(options.Path, options.Recursive, intermediate.Ext)
	if err != nil {
		return nil, err
	}

	manifest := &specs.Manifest{}

	for _, file := range files {
		reader, err := os.Open(filepath.Join(file.Path, file.Name()))
		if err != nil {
			return nil, err
		}

		definition, err := intermediate.UnmarshalHCL(file.Name(), reader)
		if err != nil {
			return nil, err
		}

		result, err := intermediate.ParseManifest(definition, options.Functions)
		if err != nil {
			return nil, err
		}

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

	err = strict.Define(options.Schema, manifest)
	if err != nil {
		return nil, err
	}

	return manifest, nil
}

// ConstructFlowManager constructs the flow managers from the given specs manifest
func ConstructFlowManager(manifest *specs.Manifest, options Options) error {
	endpoints := make([]*protocol.Endpoint, len(manifest.Endpoints))

	for index, endpoint := range manifest.Endpoints {
		f := GetFlow(manifest, endpoint.Flow)
		nodes := make([]*flow.Node, len(f.Calls))

		for index, call := range f.Calls {
			nodes[index] = flow.NewNode(call, ConstructCall(manifest, call, options), ConstructCall(manifest, call.Rollback, options))
		}

		collection, has := options.Codec[endpoint.Codec]
		if !has {
			return trace.New(trace.WithMessage("unkown endpoint codec %s", endpoint.Codec))
		}

		req, err := collection.New(specs.InputResource, f.GetInput())
		if err != nil {
			return err
		}

		res, err := collection.New(specs.InputResource, f.GetOutput())
		if err != nil {
			return err
		}

		manager := flow.NewManager(f.GetName(), nodes)

		endpoints[index] = &protocol.Endpoint{
			Flow:     manager,
			Listener: endpoint.Listener,
			Options:  endpoint.Options,
			Request:  req,
			Response: res,
		}
	}

	err := ConstructListeners(endpoints, options)
	if err != nil {
		return err
	}

	return nil
}

func ConstructCall(manifest *specs.Manifest, call specs.FlowCaller, options Options) flow.Call {
	// service := GetService(manifest, strict.GetService(call.GetEndpoint()))
	// if service == nil {
	// 	// handle err
	// }

	// caller := options.Callers[service.Caller]
	// codec := options.Codec[service.Codec]
	// call.GetDescriptor()

	// req, err := codec.New(call.GetName(), call.GetRequest())
	// res, err := codec.New(call.GetName(), call.GetResponse())

	return func(ctx context.Context, refs *refs.Store) error {
		// reader, err := req.Marshal(refs)

		// caller.Call()

		// res.Unmarshal()

		return nil
	}
}

// ConstructListeners constructs the listeners from the given collection of endpoints
func ConstructListeners(endpoints []*protocol.Endpoint, options Options) error {
	collections := make(map[string][]*protocol.Endpoint, len(options.Listeners))

	for _, endpoint := range endpoints {
		listener := GetListener(options.Listeners, endpoint.Listener)
		if listener == nil {
			return trace.New(trace.WithMessage("unkown listener %s", endpoint.Listener))
		}

		collections[endpoint.Listener] = append(collections[endpoint.Listener], endpoint)
	}

	for key, collection := range collections {
		listener := GetListener(options.Listeners, key)
		err := listener.Handle(collection)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetListener attempts to retrieve the requested listener
func GetListener(listeners []protocol.Listener, name string) protocol.Listener {
	for _, listener := range listeners {
		if listener.Name() == name {
			return listener
		}
	}

	return nil
}

// GetService attempts to retrieve a service from the given manifest matching the given name
func GetService(manifest *specs.Manifest, name string) *specs.Service {
	for _, service := range manifest.Services {
		if service.Alias == name {
			return service
		}
	}

	return nil
}

// GetFlow attempts to retrieve a flow from the given manifest matching the given name
func GetFlow(manifest *specs.Manifest, name string) *specs.Flow {
	for _, flow := range manifest.Flows {
		if flow.GetName() == name {
			return flow
		}
	}

	return nil
}
