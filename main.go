package maestro

import (
	"sync"

	"github.com/jexia/maestro/codec"
	"github.com/jexia/maestro/flow"
	"github.com/jexia/maestro/header"
	"github.com/jexia/maestro/protocol"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/strict"
	"github.com/jexia/maestro/specs/trace"
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

// ConstructFlowManager constructs the flow managers from the given specs manifest
func ConstructFlowManager(manifest *specs.Manifest, options Options) ([]*protocol.Endpoint, error) {
	endpoints := make([]*protocol.Endpoint, len(manifest.Endpoints))

	for index, endpoint := range manifest.Endpoints {
		current := manifest.GetFlow(endpoint.Flow)
		if current == nil {
			continue
		}

		nodes := make([]*flow.Node, len(current.GetNodes()))

		result := &protocol.Endpoint{
			Listener: endpoint.Listener,
			Options:  endpoint.Options,
			Request:  current.GetInput(),
			Response: current.GetOutput(),
		}

		for index, node := range current.GetNodes() {
			caller, err := ConstructCall(manifest, node, node.Call, options, current)
			if err != nil {
				return nil, err
			}

			rollback, err := ConstructCall(manifest, node, node.Rollback, options, current)
			if err != nil {
				return nil, err
			}

			nodes[index] = flow.NewNode(node, caller, rollback)
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
func ConstructCall(manifest *specs.Manifest, node *specs.Node, call *specs.Call, options Options, manager specs.FlowManager) (flow.Call, error) {
	if call == nil {
		return nil, nil
	}

	service := options.Schema.GetService(call.Service)
	if service == nil {
		return nil, trace.New(trace.WithMessage("the service for %s was not found", call.GetMethod()))
	}

	constructor := options.Callers.Get(service.GetProtocol())
	codec := options.Codec[service.GetCodec()]
	schema := options.Schema.GetService(service.GetName())

	protocol, err := constructor.New(schema, call.GetMethod(), options.Functions, service.GetOptions())
	if err != nil {
		return nil, err
	}

	request, err := ConstructRequest(node, codec, call.GetRequest())
	if err != nil {
		return nil, err
	}

	response, err := ConstructRequest(node, codec, call.GetResponse())
	if err != nil {
		return nil, err
	}

	caller := flow.NewCall(node, protocol, request, response)
	err = strict.DefineCaller(node, manifest, protocol, manager)
	if err != nil {
		return nil, err
	}

	return caller, nil
}

// ConstructRequest constructs a new request from the given parameter map and codec
func ConstructRequest(node *specs.Node, codec codec.Constructor, params *specs.ParameterMap) (*flow.Request, error) {
	message, err := codec.New(node.GetName(), params)
	if err != nil {
		return nil, err
	}

	header := header.NewManager(node.GetName(), params)
	return flow.NewRequest(message, header), nil
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
	caller, err := constructor.New(schema, "", options.Functions, service.GetOptions())
	if err != nil {
		return nil, err
	}

	return caller, nil
}

// ConstructListeners constructs the listeners from the given collection of endpoints
func ConstructListeners(endpoints []*protocol.Endpoint, options Options) error {
	collections := make(map[string][]*protocol.Endpoint, len(options.Listeners))

	for _, endpoint := range endpoints {
		if endpoint == nil {
			continue
		}

		listener := options.Listeners.Get(endpoint.Listener)
		if listener == nil {
			return trace.New(trace.WithMessage("unknown listener %s", endpoint.Listener))
		}

		collections[endpoint.Listener] = append(collections[endpoint.Listener], endpoint)
	}

	for key, collection := range collections {
		listener := options.Listeners.Get(key)
		err := listener.Handle(collection, options.Codec)
		if err != nil {
			return err
		}
	}

	return nil
}
