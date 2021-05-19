package manager

import (
	"fmt"

	"github.com/jexia/semaphore/v2/pkg/broker"
	"github.com/jexia/semaphore/v2/pkg/codec"
	"github.com/jexia/semaphore/v2/pkg/codec/metadata"
	"github.com/jexia/semaphore/v2/pkg/dependencies"
	"github.com/jexia/semaphore/v2/pkg/flow"
	"github.com/jexia/semaphore/v2/pkg/references"
	"github.com/jexia/semaphore/v2/pkg/references/forwarding"
	"github.com/jexia/semaphore/v2/pkg/specs"
	"github.com/jexia/semaphore/v2/pkg/specs/template"
)

// NewFlow constructs a new flow manager from the given configurations
func NewFlow(ctx *broker.Context, manager specs.FlowInterface, opts ...FlowOption) (*flow.Manager, error) {
	if manager == nil {
		return nil, ErrNilFlowManager
	}

	options := NewFlowOptions(opts...)
	nodes := make([]*flow.Node, len(manager.GetNodes()))

	for index, spec := range manager.GetNodes() {
		result, err := node(ctx, manager, spec, options)
		if err != nil {
			return nil, err
		}

		nodes[index] = result
	}

	stack := options.stack.Load(manager.GetOutput())
	flow := flow.NewManager(ctx, manager.GetName(), nodes, manager.GetOnError(), stack, &flow.ManagerMiddleware{
		BeforeDo:       options.BeforeManagerDo,
		AfterDo:        options.AfterManagerDo,
		BeforeRollback: options.BeforeManagerRollback,
		AfterRollback:  options.AfterManagerRollback,
	})

	if options.AfterFlowConstruction != nil {
		err := options.AfterFlowConstruction(ctx, manager, flow)
		if err != nil {
			return nil, err
		}
	}

	return flow, nil
}

func node(ctx *broker.Context, manager specs.FlowInterface, node *specs.Node, options FlowOptions) (*flow.Node, error) {
	arguments := flow.NodeArguments{
		flow.WithNodeMiddleware(flow.NodeMiddleware{
			BeforeDo:       options.BeforeNodeDo,
			AfterDo:        options.AfterNodeDo,
			BeforeRollback: options.BeforeNodeRollback,
			AfterRollback:  options.AfterNodeRollback,
		}),
	}

	if node.Intermediate != nil {
		stack := options.stack.Load(node.Intermediate)
		arguments.Set(flow.WithFunctions(stack))
	}

	if node.Condition != nil {
		arguments.Set(flow.WithCondition(
			flow.NewCondition(options.stack.Load(node.Condition.Params),
				node.Condition,
			),
		))
	}

	if node.Call != nil {
		caller, err := service(ctx, manager, node, node.Call, options)
		if err != nil {
			return nil, err
		}

		arguments.Set(flow.WithCall(caller))
	}

	if node.Rollback != nil {
		rollback, err := service(ctx, manager, node, node.Rollback, options)
		if err != nil {
			return nil, err
		}

		arguments.Set(flow.WithRollback(rollback))
	}

	return flow.NewNode(ctx, node, arguments...), nil
}

// service constructs a new flow caller for the given service
func service(ctx *broker.Context, manager specs.FlowInterface, node *specs.Node, call *specs.Call, options FlowOptions) (flow.Call, error) {
	if call == nil {
		return nil, nil
	}

	if call.Service == "" {
		return nil, ErrNoServiceName{
			Flow: manager.GetName(),
			Node: node.ID,
		}
	}

	service := options.services.Get(call.Service)
	if service == nil {
		return nil, ErrNoService{
			Flow:    manager.GetName(),
			Node:    node.ID,
			Service: call.Service,
		}
	}

	constructor := options.Callers.Get(service.Transport)

	if constructor == nil {
		return nil, ErrNoTransport{
			Transport: service.Transport,
			Details: Details{
				Service: call.Service,
				Flow:    manager.GetName(),
				Node:    node.ID,
			},
		}
	}

	client, ok := options.discoveries[service.Resolver]
	if !ok {
		return nil, fmt.Errorf("service discovery with name '%s' is not registered", service.Resolver)
	}

	resolver, err := client.Resolver(service.Host)
	if err != nil {
		return nil, fmt.Errorf("failed to create service resolver: %w", err)
	}

	dialer, err := constructor.Dial(service, options.Functions, service.Options, resolver)
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

			forwarding.ResolvePropertyReferences(&reference.Template, node.DependsOn)
			err = dependencies.Resolve(manager, node.DependsOn, node.ID, make(dependencies.Unresolved))
			if err != nil {
				return nil, err
			}
		}
	}

	reqcodec := options.Codec.Get(service.RequestCodec)
	if reqcodec == nil {
		return nil, ErrNoRequestCodec{
			Codec: service.RequestCodec,
			Details: Details{
				Service: call.Service,
				Flow:    manager.GetName(),
				Node:    node.ID,
			},
		}
	}

	rescodec := options.Codec.Get(service.ResponseCodec)
	if rescodec == nil {
		return nil, ErrNoResponseCodec{
			Codec: service.RequestCodec,
			Details: Details{
				Service: call.Service,
				Flow:    manager.GetName(),
				Node:    node.ID,
			},
		}
	}

	unexpected, err := errorHandler(ctx, node, rescodec, node.OnError, options)
	if err != nil {
		return nil, err
	}

	request, err := messageHandle(ctx, node, reqcodec, call.Request, options)
	if err != nil {
		return nil, err
	}

	response, err := messageHandle(ctx, node, rescodec, call.Response, options)
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

// messageHandle constructs a new request from the given parameter map and codec
func messageHandle(ctx *broker.Context, node *specs.Node, constructor codec.Constructor, params *specs.ParameterMap, options FlowOptions) (*flow.Request, error) {
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

	metadata := metadata.NewManager(ctx, node.ID, params.Header)
	request := &flow.Request{
		Functions: options.stack.Load(params),
		Codec:     codec,
		Metadata:  metadata,
	}

	return request, nil
}

// errorHandler constructs a new error object from the given parameter map and codec
func errorHandler(ctx *broker.Context, node *specs.Node, constructor codec.Constructor, handle *specs.OnError, options FlowOptions) (*flow.OnError, error) {
	if handle == nil {
		return nil, nil
	}

	var codec codec.Manager
	var meta *metadata.Manager

	if handle.Response != nil && constructor != nil {
		manager, err := constructor.New(template.JoinPath(node.ID, template.ErrorResource), handle.Response)
		if err != nil {
			return nil, err
		}

		codec = manager
		meta = metadata.NewManager(ctx, node.ID, handle.Response.Header)
	}

	return flow.NewOnError(options.stack.Load(handle.Response), codec, meta, handle), nil
}
