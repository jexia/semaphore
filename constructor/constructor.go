package constructor

import (
	"github.com/jexia/maestro/codec"
	"github.com/jexia/maestro/flow"
	"github.com/jexia/maestro/metadata"
	"github.com/jexia/maestro/protocol"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/strict"
	"github.com/jexia/maestro/specs/trace"
)

// Specs construct a specs manifest from the given options
func Specs(options Options) (*specs.Manifest, error) {
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

// FlowManager constructs the flow managers from the given specs manifest
func FlowManager(manifest *specs.Manifest, options Options) ([]*protocol.Endpoint, error) {
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
			caller, err := Call(manifest, node, node.Call, options, current)
			if err != nil {
				return nil, err
			}

			rollback, err := Call(manifest, node, node.Rollback, options, current)
			if err != nil {
				return nil, err
			}

			nodes[index] = flow.NewNode(node, caller, rollback)
		}

		forward, err := Forward(manifest, current.GetForward(), options)
		if err != nil {
			return nil, err
		}

		result.Forward = forward
		result.Flow = flow.NewManager(current.GetName(), nodes)

		endpoints[index] = result
	}

	err := Listeners(endpoints, options)
	if err != nil {
		return nil, err
	}

	return endpoints, nil
}

// Call constructs a flow caller for the given node call.
func Call(manifest *specs.Manifest, node *specs.Node, call *specs.Call, options Options, manager specs.FlowManager) (flow.Call, error) {
	if call == nil {
		return nil, nil
	}

	service := options.Schema.GetService(call.Service)
	if service == nil {
		return nil, trace.New(trace.WithMessage("the service for %s was not found", call.GetMethod()))
	}

	constructor := options.Callers.Get(service.GetProtocol())
	codec := options.Codec[service.GetCodec()]
	schema := options.Schema.GetService(service.GetFullyQualifiedName())

	if schema == nil {
		return nil, trace.New(trace.WithMessage("service not found '%s'", service.GetFullyQualifiedName()))
	}

	if constructor == nil {
		return nil, trace.New(trace.WithMessage("protocol constructor not found '%s' for service '%s'", service.GetProtocol(), service.GetName()))
	}

	protocol, err := constructor.Dial(schema, options.Functions, service.GetOptions())
	if err != nil {
		return nil, err
	}

	request, err := Request(node, codec, call.GetRequest())
	if err != nil {
		return nil, err
	}

	response, err := Request(node, codec, call.GetResponse())
	if err != nil {
		return nil, err
	}

	caller := flow.NewCall(node, protocol, call.Method, request, response)
	err = strict.DefineCaller(node, manifest, protocol, manager)
	if err != nil {
		return nil, err
	}

	return caller, nil
}

// Request constructs a new request from the given parameter map and codec
func Request(node *specs.Node, codec codec.Constructor, params *specs.ParameterMap) (*flow.Request, error) {
	message, err := codec.New(node.GetName(), params)
	if err != nil {
		return nil, err
	}

	metadata := metadata.NewManager(node.GetName(), params)
	return flow.NewRequest(message, metadata), nil
}

// Forward constructs a flow caller for the given call.
func Forward(manifest *specs.Manifest, call *specs.Call, options Options) (protocol.Call, error) {
	if call == nil {
		return nil, nil
	}

	service := options.Schema.GetService(call.GetService())
	if service == nil {
		return nil, trace.New(trace.WithMessage("the service for %s was not found", call.GetMethod()))
	}

	schema := options.Schema.GetService(service.GetName())
	constructor := options.Callers.Get(service.GetProtocol())
	caller, err := constructor.Dial(schema, options.Functions, service.GetOptions())
	if err != nil {
		return nil, err
	}

	return caller, nil
}

// Listeners constructs the listeners from the given collection of endpoints
func Listeners(endpoints []*protocol.Endpoint, options Options) error {
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
