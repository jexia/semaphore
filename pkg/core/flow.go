package core

import (
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/config"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/codec"
	"github.com/jexia/semaphore/pkg/codec/metadata"
	"github.com/jexia/semaphore/pkg/core/flows/condition"
	"github.com/jexia/semaphore/pkg/core/flows/listeners"
	"github.com/jexia/semaphore/pkg/core/trace"
	"github.com/jexia/semaphore/pkg/dependencies"
	"github.com/jexia/semaphore/pkg/flow"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/references/forwarding"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
	"github.com/jexia/semaphore/pkg/transport"
)

// Apply constructs the flow managers from the given specs manifest
func Apply(ctx *broker.Context, mem functions.Collection, services specs.ServiceList, endpoints specs.EndpointList, flows specs.FlowListInterface, options config.Options) ([]*transport.Endpoint, error) {
	results := make([]*transport.Endpoint, len(endpoints))

	logger.Debug(ctx, "constructing endpoints")

	for index, endpoint := range endpoints {
		manager := flows.Get(endpoint.Flow)
		if manager == nil {
			continue
		}

		nodes := make([]*flow.Node, len(manager.GetNodes()))

		for index, node := range manager.GetNodes() {
			arguments := []flow.NodeOption{
				flow.WithNodeMiddleware(&flow.NodeMiddleware{
					BeforeDo:       options.BeforeNodeDo,
					AfterDo:        options.AfterNodeDo,
					BeforeRollback: options.BeforeNodeRollback,
					AfterRollback:  options.AfterNodeRollback,
				}),
			}

			if node.Intermediate != nil {
				stack := mem[node.Intermediate]
				arguments = append(arguments, flow.WithFunctions(stack))
			}

			if node.Condition != nil {
				arguments = append(arguments, flow.WithCondition(condition.New(ctx, mem, node.Condition)))
			}

			if node.Call != nil {
				caller, err := NewServiceCall(ctx, mem, services, node, node.Call, options, manager)
				if err != nil {
					return nil, err
				}

				arguments = append(arguments, flow.WithCall(caller))
			}

			if node.Rollback != nil {
				rollback, err := NewServiceCall(ctx, mem, services, node, node.Rollback, options, manager)
				if err != nil {
					return nil, err
				}

				arguments = append(arguments, flow.WithRollback(rollback))
			}

			nodes[index] = flow.NewNode(ctx, node, arguments...)
		}

		forward, err := NewForward(services, manager.GetForward(), options)
		if err != nil {
			return nil, err
		}

		stack := mem[manager.GetOutput()]
		flow := flow.NewManager(ctx, manager.GetName(), nodes, manager.GetOnError(), stack, &flow.ManagerMiddleware{
			BeforeDo:       options.BeforeManagerDo,
			AfterDo:        options.AfterManagerDo,
			BeforeRollback: options.BeforeManagerRollback,
			AfterRollback:  options.AfterManagerRollback,
		})

		results[index] = transport.NewEndpoint(endpoint.Listener, flow, forward, endpoint.Options, manager.GetInput(), manager.GetOutput())
	}

	err := listeners.Apply(results, options)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// NewServiceCall constructs a new flow caller for the given service
func NewServiceCall(ctx *broker.Context, mem functions.Collection, services specs.ServiceList, node *specs.Node, call *specs.Call, options config.Options, manager specs.FlowInterface) (flow.Call, error) {
	if call == nil {
		return nil, nil
	}

	if call.Service == "" {
		return nil, trace.New(trace.WithMessage("invalid service name, no service name configured in '%s'", node.ID))
	}

	service := services.Get(call.Service)
	if service == nil {
		return nil, trace.New(trace.WithMessage("the service for '%s' was not found in '%s'", call.Service, node.ID))
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
			err := references.ResolveProperty(ctx, node, reference, manager)
			if err != nil {
				return nil, err
			}

			forwarding.ResolvePropertyReferences(reference, node.DependsOn)
			err = dependencies.ResolveNode(manager, node, make(specs.Dependencies))
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
func NewRequest(ctx *broker.Context, node *specs.Node, mem functions.Collection, constructor codec.Constructor, params *specs.ParameterMap) (*flow.Request, error) {
	if params == nil {
		return nil, nil
	}

	var codec codec.Manager
	if constructor != nil {
		manager, err := constructor.New(node.ID, params)
		if err != nil {
			return nil, err
		}

		codec = manager
	}

	stack := mem[params]
	metadata := metadata.NewManager(ctx, node.ID, params.Header)
	return flow.NewRequest(stack, codec, metadata), nil
}

// NewForward constructs a flow caller for the given call.
func NewForward(services specs.ServiceList, call *specs.Call, options config.Options) (*transport.Forward, error) {
	if call == nil {
		return nil, nil
	}

	service := services.Get(call.Service)
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

// NewError constructs a new error object from the given parameter map and codec
func NewError(ctx *broker.Context, node *specs.Node, mem functions.Collection, constructor codec.Constructor, err *specs.OnError) (*flow.OnError, error) {
	if err == nil {
		return nil, nil
	}

	var codec codec.Manager
	var meta *metadata.Manager
	var stack functions.Stack

	if err.Response != nil && constructor != nil {
		params := err.Response

		// TODO: check if I would like props to be defined like this
		manager, err := constructor.New(template.JoinPath(node.ID, template.ErrorResource), params)
		if err != nil {
			return nil, err
		}

		codec = manager
		stack = mem[params]
		meta = metadata.NewManager(ctx, node.ID, params.Header)
	}

	return flow.NewOnError(stack, codec, meta, err.Status, err.Message), nil
}
