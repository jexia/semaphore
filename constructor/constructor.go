package constructor

import (
	"github.com/jexia/maestro/codec"
	"github.com/jexia/maestro/flow"
	"github.com/jexia/maestro/instance"
	"github.com/jexia/maestro/metadata"
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/strict"
	"github.com/jexia/maestro/specs/trace"
	"github.com/jexia/maestro/transport"
)

// Specs construct a specs manifest from the given options
func Specs(ctx instance.Context, options Options) (*specs.Manifest, error) {
	result := &specs.Manifest{}

	for _, resolver := range options.Definitions {
		if resolver == nil {
			continue
		}

		manifest, err := resolver(ctx, options.Functions)
		if err != nil {
			return nil, err
		}

		result.Merge(manifest)
	}

	for _, resolver := range options.Schemas {
		if resolver == nil {
			continue
		}

		err := resolver(ctx, options.Schema)
		if err != nil {
			return nil, err
		}
	}

	err := specs.CheckManifestDuplicates(ctx, result)
	if err != nil {
		return nil, err
	}

	err = specs.ResolveManifestDependencies(ctx, result)
	if err != nil {
		return nil, err
	}

	err = strict.DefineManifest(ctx, options.Schema, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// FlowManager constructs the flow managers from the given specs manifest
func FlowManager(ctx instance.Context, manifest *specs.Manifest, options Options) ([]*transport.Endpoint, error) {
	endpoints := make([]*transport.Endpoint, len(manifest.Endpoints))

	for index, endpoint := range manifest.Endpoints {
		manager := manifest.GetFlow(endpoint.Flow)
		if manager == nil {
			continue
		}

		nodes := make([]*flow.Node, len(manager.GetNodes()))

		result := &transport.Endpoint{
			Listener: endpoint.Listener,
			Options:  endpoint.Options,
			Request:  manager.GetInput(),
			Response: manager.GetOutput(),
		}

		for index, node := range manager.GetNodes() {
			caller, err := Call(ctx, manifest, node, node.Call, options, manager)
			if err != nil {
				return nil, err
			}

			rollback, err := Call(ctx, manifest, node, node.Rollback, options, manager)
			if err != nil {
				return nil, err
			}

			nodes[index] = flow.NewNode(ctx, node, caller, rollback)
		}

		forward, err := Forward(manifest, manager.GetForward(), options)
		if err != nil {
			return nil, err
		}

		result.Flow = flow.NewManager(ctx, manager.GetName(), nodes)
		result.Forward = forward

		endpoints[index] = result
	}

	err := Listeners(endpoints, options)
	if err != nil {
		return nil, err
	}

	return endpoints, nil
}

// Call constructs a flow caller for the given node call.
func Call(ctx instance.Context, manifest *specs.Manifest, node *specs.Node, call *specs.Call, options Options, manager specs.FlowManager) (flow.Call, error) {
	if call == nil {
		return nil, nil
	}

	service := options.Schema.GetService(call.Service)
	if service == nil {
		return nil, trace.New(trace.WithMessage("the service for %s was not found", call.GetMethod()))
	}

	constructor := options.Callers.Get(service.GetTransport())
	codec := options.Codec[service.GetCodec()]
	schema := options.Schema.GetService(service.GetFullyQualifiedName())

	if schema == nil {
		return nil, trace.New(trace.WithMessage("service not found '%s'", service.GetFullyQualifiedName()))
	}

	if constructor == nil {
		return nil, trace.New(trace.WithMessage("transport constructor not found '%s' for service '%s'", service.GetTransport(), service.GetName()))
	}

	transport, err := constructor.Dial(schema, options.Functions, service.GetOptions())
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

	caller := flow.NewCall(ctx, node, transport, call.Method, request, response)
	err = strict.DefineCaller(ctx, node, manifest, transport, manager)
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
func Forward(manifest *specs.Manifest, call *specs.Call, options Options) (schema.Service, error) {
	if call == nil {
		return nil, nil
	}

	service := options.Schema.GetService(call.GetService())
	if service == nil {
		return nil, trace.New(trace.WithMessage("the service for %s was not found", call.GetMethod()))
	}

	return service, nil
}

// Listeners constructs the listeners from the given collection of endpoints
func Listeners(endpoints []*transport.Endpoint, options Options) error {
	collections := make(map[string][]*transport.Endpoint, len(options.Listeners))

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
