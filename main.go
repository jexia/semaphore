package maestro

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/jexia/maestro/codec"
	"github.com/jexia/maestro/definitions/hcl"
	"github.com/jexia/maestro/flow"
	"github.com/jexia/maestro/protocol"
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/strict"
	"github.com/jexia/maestro/specs/trace"
	"github.com/jexia/maestro/utils"
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

// Option represents a constructor func which sets a given option
type Option func(*Options)

// Options represents all the available options
type Options struct {
	Path      string
	Recursive bool
	Codec     map[string]codec.Constructor
	Callers   protocol.Callers
	Listeners protocol.Listeners
	Schema    *schema.Store
	Functions specs.CustomDefinedFunctions
}

// NewOptions constructs a options object from the given option constructors
func NewOptions(options ...Option) Options {
	result := Options{
		Codec:  make(map[string]codec.Constructor),
		Schema: schema.NewStore(),
	}

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

// WithCaller appends the given caller to the collection of available callers
func WithCaller(caller protocol.Caller) Option {
	return func(options *Options) {
		options.Callers = append(options.Callers, caller)
	}
}

// WithListener appends the given listener to the collection of available listeners
func WithListener(listener protocol.Listener) Option {
	return func(options *Options) {
		options.Listeners = append(options.Listeners, listener)
	}
}

// WithSchema appends the schema collection to the schema store
func WithSchema(collection schema.Collection) Option {
	return func(options *Options) {
		options.Schema.Add(collection)
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

	err = strict.Define(options.Schema, manifest)
	if err != nil {
		return nil, err
	}

	return manifest, nil
}

// ConstructFlowManager constructs the flow managers from the given specs manifest
func ConstructFlowManager(manifest *specs.Manifest, options Options) ([]*protocol.Endpoint, error) {
	endpoints := make([]*protocol.Endpoint, len(manifest.Endpoints))

	for index, endpoint := range manifest.Endpoints {
		current := manifest.GetFlow(endpoint.Flow)
		nodes := make([]*flow.Node, len(current.GetNodes()))

		result := &protocol.Endpoint{
			Listener: endpoint.Listener,
			Options:  endpoint.Options,
		}

		for index, node := range current.GetNodes() {
			caller, err := ConstructCall(manifest, node, node.Call, options)
			if err != nil {
				return nil, err
			}

			rollback, err := ConstructCall(manifest, node, node.Rollback, options)
			if err != nil {
				return nil, err
			}

			nodes[index] = flow.NewNode(node, caller, rollback)
		}

		collection, has := options.Codec[endpoint.Codec]
		if !has {
			return nil, trace.New(trace.WithMessage("unknown endpoint codec %s", endpoint.Codec))
		}

		if current.GetInput() != nil {
			req, err := collection.New(specs.InputResource, current.GetInput())
			if err != nil {
				return nil, err
			}

			result.Request = req
		}

		if current.GetOutput() != nil {
			res, err := collection.New(specs.InputResource, current.GetOutput())
			if err != nil {
				return nil, err
			}

			result.Response = res
			result.Header = protocol.NewHeaderManager(specs.InputResource, current.GetOutput())
		}

		forward, err := ConstructForward(manifest, current.GetForward(), options)
		if err != nil {
			return nil, err
		}

		result.Forward = forward
		result.Flow = flow.NewManager(current.GetName(), nodes)

		endpoints[index] = result
	}

	err := ConstructListeners(endpoints, options)
	if err != nil {
		return nil, err
	}

	return endpoints, nil
}

// ConstructCall constructs a flow caller for the given node call.
func ConstructCall(manifest *specs.Manifest, node *specs.Node, call *specs.Call, options Options) (flow.Call, error) {
	if call == nil {
		return nil, nil
	}

	service := options.Schema.GetService(call.Service)
	if service == nil {
		return nil, trace.New(trace.WithMessage("the service for %s was not found", call.GetMethod()))
	}

	constructor := options.Callers.Get(service.GetProtocol())
	codec := options.Codec[service.GetCodec()]

	req, err := codec.New(node.GetName(), call.GetRequest())
	if err != nil {
		return nil, err
	}

	res, err := codec.New(node.GetName(), call.GetResponse())
	if err != nil {
		return nil, err
	}

	header := protocol.NewHeaderManager(node.GetName(), call.GetRequest())
	schema := options.Schema.GetService(service.GetName())

	caller, err := constructor.New(schema, call.GetMethod(), service.GetOptions())
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context, refs *refs.Store) error {
		body, err := req.Marshal(refs)
		if err != nil {
			return err
		}

		reader, writer := io.Pipe()
		w := protocol.NewResponseWriter(writer)
		r := &protocol.Request{
			Context: ctx,
			Body:    body,
			Header:  header.Marshal(refs),
		}

		go func() {
			defer writer.Close()
			err := caller.Call(w, r, refs)
			if err != nil {
				log.Println(err)
			}
		}()

		err = res.Unmarshal(reader, refs)
		if err != nil {
			return err
		}

		if !protocol.StatusSuccess(w.Status()) {
			log.WithFields(log.Fields{
				"node":   node.GetName(),
				"status": w.Status(),
			}).Error("Faulty status code")

			return errors.New("rollback required")
		}

		return nil
	}, nil
}

// ConstructForward constructs a flow caller for the given call.
func ConstructForward(manifest *specs.Manifest, call *specs.Call, options Options) (protocol.Call, error) {
	if call == nil {
		return nil, nil
	}

	service := options.Schema.GetService(call.GetService())
	if service == nil {
		return nil, trace.New(trace.WithMessage("the service for %s was not found", call.GetMethod()))
	}

	schema := options.Schema.GetService(service.GetName())
	constructor := options.Callers.Get(service.GetProtocol())
	caller, err := constructor.New(schema, call.GetMethod(), service.GetOptions()) // MARK
	if err != nil {
		return nil, err
	}

	return caller, nil
}

// ConstructListeners constructs the listeners from the given collection of endpoints
func ConstructListeners(endpoints []*protocol.Endpoint, options Options) error {
	collections := make(map[string][]*protocol.Endpoint, len(options.Listeners))

	for _, endpoint := range endpoints {
		listener := options.Listeners.Get(endpoint.Listener)
		if listener == nil {
			return trace.New(trace.WithMessage("unknown listener %s", endpoint.Listener))
		}

		collections[endpoint.Listener] = append(collections[endpoint.Listener], endpoint)
	}

	for key, collection := range collections {
		listener := options.Listeners.Get(key)
		err := listener.Handle(collection)
		if err != nil {
			return err
		}
	}

	return nil
}
