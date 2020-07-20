package core

import (
	"github.com/jexia/semaphore/pkg/codec"
	"github.com/jexia/semaphore/pkg/core/api"
	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/core/trace"
	"github.com/jexia/semaphore/pkg/dependencies"
	"github.com/jexia/semaphore/pkg/flow"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/metadata"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/transport"
)

// NewServiceCall constructs a new flow caller for the given service
func NewServiceCall(ctx instance.Context, mem functions.Collection, services *specs.ServicesManifest, flows *specs.FlowsManifest, node *specs.Node, call *specs.Call, options api.Options, manager specs.FlowResourceManager) (flow.Call, error) {
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

	method := dialer.GetMethod(node.Call.Method)
	if method != nil {
		for _, reference := range method.References() {
			err := references.DefineProperty(ctx, node, reference, manager)
			if err != nil {
				return nil, err
			}

			dependencies.ResolvePropertyReferences(reference, node.DependsOn)
			err = dependencies.ResolveNode(manager, node, make(map[string]*specs.Node))
			if err != nil {
				return nil, err
			}
		}
	}

	codec := options.Codec.Get(service.Codec)
	if codec == nil {
		return nil, trace.New(trace.WithMessage("codec not found '%s'", service.Codec))
	}

	unexpected, err := NewError(ctx, node, mem, codec, node.OnError)
	if err != nil {
		return nil, err
	}

	request, err := NewRequest(ctx, node, mem, codec, call.Request)
	if err != nil {
		return nil, err
	}

	response, err := NewRequest(ctx, node, mem, codec, call.Response)
	if err != nil {
		return nil, err
	}

	caller := flow.NewCall(ctx, node, &flow.CallOptions{
		ExpectedStatus: node.ExpectStatus,
		Transport:      dialer,
		Method:         dialer.GetMethod(call.Method),
		Err:            unexpected,
		Request:        request,
		Response:       response,
	})

	return caller, nil
}

// NewRequest constructs a new request from the given parameter map and codec
func NewRequest(ctx instance.Context, node *specs.Node, mem functions.Collection, constructor codec.Constructor, params *specs.ParameterMap) (*flow.Request, error) {
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
	metadata := metadata.NewManager(ctx, node.Name, params.Header)
	return flow.NewRequest(stack, codec, metadata), nil
}

// NewForward constructs a flow caller for the given call.
func NewForward(services *specs.ServicesManifest, flows *specs.FlowsManifest, call *specs.Call, options api.Options) (*transport.Forward, error) {
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
		result.Schema = call.Request.Header
	}

	return result, nil
}
