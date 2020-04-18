package constructor

import (
	"github.com/jexia/maestro/pkg/codec"
	"github.com/jexia/maestro/pkg/flow"
	"github.com/jexia/maestro/pkg/functions"
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/metadata"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/checks"
	"github.com/jexia/maestro/pkg/specs/dependencies"
	"github.com/jexia/maestro/pkg/specs/references"
	"github.com/jexia/maestro/pkg/specs/trace"
	"github.com/jexia/maestro/pkg/transport"
)

// Specs construct a specs manifest from the given options
func Specs(ctx instance.Context, options Options) (functions.Collection, *specs.FlowsManifest, *specs.ServicesManifest, *specs.SchemaManifest, error) {
	flows := specs.NewFlowsManifest()
	services := specs.NewServicesManifest()
	schema := specs.NewSchemaManifest()

	for _, resolver := range options.Flows {
		if resolver == nil {
			continue
		}

		manifest, err := resolver(ctx)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		flows.Merge(manifest)
	}

	for _, resolver := range options.Services {
		if resolver == nil {
			continue
		}

		manifest, err := resolver(ctx)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		services.Merge(manifest)
	}

	for _, resolver := range options.Schemas {
		if resolver == nil {
			continue
		}

		manifest, err := resolver(ctx)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		schema.Merge(manifest)
	}

	err := references.DefineManifest(ctx, services, schema, flows)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	err = checks.ManifestDuplicates(ctx, flows)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	mem := functions.Collection{}
	err = functions.PrepareManifestFunctions(ctx, mem, options.Functions, flows)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	err = references.CompareManifestTypes(ctx, services, schema, flows)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	dependencies.ResolveReferences(ctx, flows)
	err = dependencies.ResolveManifest(ctx, flows)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return mem, flows, services, schema, nil
}

// FlowManager constructs the flow managers from the given specs manifest
func FlowManager(ctx instance.Context, mem functions.Collection, services *specs.ServicesManifest, flows *specs.FlowsManifest, options Options) ([]*transport.Endpoint, error) {
	endpoints := make([]*transport.Endpoint, len(flows.Endpoints))

	for index, endpoint := range flows.Endpoints {
		manager := flows.GetFlow(endpoint.Flow)
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
			caller, err := Call(ctx, mem, services, flows, node, node.Call, options, manager)
			if err != nil {
				return nil, err
			}

			rollback, err := Call(ctx, mem, services, flows, node, node.Rollback, options, manager)
			if err != nil {
				return nil, err
			}

			nodes[index] = flow.NewNode(ctx, node, caller, rollback)
		}

		forward, err := Forward(services, flows, manager.GetForward(), options)
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
func Call(ctx instance.Context, mem functions.Collection, services *specs.ServicesManifest, flows *specs.FlowsManifest, node *specs.Node, call *specs.Call, options Options, manager specs.FlowResourceManager) (flow.Call, error) {
	if call == nil {
		return nil, nil
	}

	if call.Service != "" {
		return NewServiceCall(ctx, mem, services, flows, node, call, options, manager)
	}

	request, err := Request(node, mem, nil, call.Request)
	if err != nil {
		return nil, err
	}

	response, err := Request(node, mem, nil, call.Response)
	if err != nil {
		return nil, err
	}

	caller := flow.NewCall(ctx, node, &flow.CallOptions{
		Request:  request,
		Response: response,
	})

	return caller, nil
}

// NewServiceCall constructs a new flow caller for the given service
func NewServiceCall(ctx instance.Context, mem functions.Collection, services *specs.ServicesManifest, flows *specs.FlowsManifest, node *specs.Node, call *specs.Call, options Options, manager specs.FlowResourceManager) (flow.Call, error) {
	if call == nil {
		return nil, nil
	}

	if call.Service == "" {
		return nil, trace.New(trace.WithMessage("invalid service name, no service name configured in '%s'", node.Name))
	}

	service := services.GetService(call.Service)
	if service == nil {
		return nil, trace.New(trace.WithMessage("the service for '%s' was not found in '%s'", call.Service, node.Name))
	}

	constructor := options.Callers.Get(service.Transport)

	if constructor == nil {
		return nil, trace.New(trace.WithMessage("transport constructor not found '%s' for service '%s'", service.Transport, service.Name))
	}

	dialer, err := constructor.Dial(service, options.Functions, service.Options)
	if err != nil {
		return nil, err
	}

	err = transport.DefineCaller(ctx, node, flows, dialer, manager)
	if err != nil {
		return nil, err
	}

	codec := options.Codec.Get(service.Codec)
	if codec == nil {
		return nil, trace.New(trace.WithMessage("codec not found '%s'", service.Codec))
	}

	request, err := Request(node, mem, codec, call.Request)
	if err != nil {
		return nil, err
	}

	response, err := Request(node, mem, codec, call.Response)
	if err != nil {
		return nil, err
	}

	caller := flow.NewCall(ctx, node, &flow.CallOptions{
		Transport: dialer,
		Method:    dialer.GetMethod(call.Method),
		Request:   request,
		Response:  response,
	})

	return caller, nil
}

// Request constructs a new request from the given parameter map and codec
func Request(node *specs.Node, mem functions.Collection, constructor codec.Constructor, params *specs.ParameterMap) (*flow.Request, error) {
	if params == nil {
		return nil, nil
	}

	var codec codec.Manager
	if constructor != nil {
		manager, err := constructor.New(node.Name, params)
		if err != nil {
			return nil, err
		}

		codec = manager
	}

	stack := mem[params]
	metadata := metadata.NewManager(node.Name, params.Header)
	return flow.NewRequest(stack, codec, metadata), nil
}

// Forward constructs a flow caller for the given call.
func Forward(services *specs.ServicesManifest, flows *specs.FlowsManifest, call *specs.Call, options Options) (*transport.Forward, error) {
	if call == nil {
		return nil, nil
	}

	service := services.GetService(call.Service)
	if service == nil {
		return nil, trace.New(trace.WithMessage("the service for '%s' was not found", call.Method))
	}

	result := &transport.Forward{
		Service: service,
	}

	if call.Request != nil {
		result.Header = call.Request.Header
	}

	return result, nil
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
