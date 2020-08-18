package endpoints

import (
	"errors"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/config"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/broker/trace"
	"github.com/jexia/semaphore/pkg/codec"
	"github.com/jexia/semaphore/pkg/codec/metadata"
	"github.com/jexia/semaphore/pkg/dependencies"
	"github.com/jexia/semaphore/pkg/flow"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/references/forwarding"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
	"github.com/jexia/semaphore/pkg/transport"
)

// ErrFlowNotFound is returned when a given flow is not found
var ErrFlowNotFound = errors.New("flow not found")

// NewOptions constructs a new endpoint option object from the passed options
func NewOptions(opts ...EndpointOption) Options {
	result := Options{}

	for _, opt := range opts {
		opt(&result)
	}

	return result
}

// Options represents a collection of options which are used
// construct and endpoints.
type Options struct {
	config.Options
	stack    functions.Collection
	services specs.ServiceList
	flows    specs.FlowListInterface
}

// EndpointOption applies the given options to the apply options object.
type EndpointOption func(*Options)

// WithFunctions sets the given functions
func WithFunctions(stack functions.Collection) EndpointOption {
	return func(options *Options) {
		options.stack = stack
	}
}

// WithServices sets the given services
func WithServices(services specs.ServiceList) EndpointOption {
	return func(options *Options) {
		options.services = services
	}
}

// WithOptions sets the given options
func WithOptions(conf config.Options) EndpointOption {
	return func(options *Options) {
		options.Options = conf
	}
}

// Transporters constructs a new transport Endpoints list from the given endpoints and options
func Transporters(ctx *broker.Context, endpoints specs.EndpointList, flows specs.FlowListInterface, opts ...EndpointOption) (transport.EndpointList, error) {
	options := NewOptions(opts...)

	results := make(transport.EndpointList, len(endpoints))
	logger.Debug(ctx, "constructing endpoints")

	for index, endpoint := range endpoints {
		result, err := transporter(ctx, endpoint, flows, options)
		if err != nil {
			return nil, err
		}

		results[index] = result
	}

	return results, nil
}

// transporter constructs a transport endpoint from the given specifications, flows, and options
func transporter(ctx *broker.Context, endpoint *specs.Endpoint, flows specs.FlowListInterface, options Options) (*transport.Endpoint, error) {
	manager := flows.Get(endpoint.Flow)
	if manager == nil {
		return nil, ErrFlowNotFound
	}

	nodes := make([]*flow.Node, len(manager.GetNodes()))

	for index, spec := range manager.GetNodes() {
		result, err := node(ctx, manager, spec, options)
		if err != nil {
			return nil, err
		}

		nodes[index] = result
	}

	forward, err := forwarder(manager.GetForward(), options)
	if err != nil {
		return nil, err
	}

	stack := options.stack.Load(manager.GetOutput())
	flow := flow.NewManager(ctx, manager.GetName(), nodes, manager.GetOnError(), stack, &flow.ManagerMiddleware{
		BeforeDo:       options.BeforeManagerDo,
		AfterDo:        options.AfterManagerDo,
		BeforeRollback: options.BeforeManagerRollback,
		AfterRollback:  options.AfterManagerRollback,
	})

	return transport.NewEndpoint(endpoint.Listener, flow, forward, endpoint.Options, manager.GetInput(), manager.GetOutput()), nil
}

func node(ctx *broker.Context, manager specs.FlowInterface, node *specs.Node, options Options) (*flow.Node, error) {
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
func service(ctx *broker.Context, manager specs.FlowInterface, node *specs.Node, call *specs.Call, options Options) (flow.Call, error) {
	if call == nil {
		return nil, nil
	}

	if call.Service == "" {
		return nil, trace.New(trace.WithMessage("invalid service name, no service name configured in '%s'", node.ID))
	}

	service := options.services.Get(call.Service)
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
			err = dependencies.Resolve(manager, node.DependsOn, node.ID, make(dependencies.Unresolved))
			if err != nil {
				return nil, err
			}
		}
	}

	reqcodec := options.Codec.Get(service.RequestCodec)
	if reqcodec == nil {
		return nil, trace.New(trace.WithMessage("response codec not found '%s'", service.ResponseCodec))
	}

	request, err := messageHandle(ctx, node, reqcodec, call.Request, options)
	if err != nil {
		return nil, err
	}

	rescodec := options.Codec.Get(service.ResponseCodec)
	if rescodec == nil {
		return nil, trace.New(trace.WithMessage("response codec not found '%s'", service.ResponseCodec))
	}

	unexpected, err := errorHandler(ctx, node, rescodec, node.OnError, options)
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
func messageHandle(ctx *broker.Context, node *specs.Node, constructor codec.Constructor, params *specs.ParameterMap, options Options) (*flow.Request, error) {
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

// newForward constructs a flow caller for the given call.
func forwarder(call *specs.Call, options Options) (*transport.Forward, error) {
	if call == nil {
		return nil, nil
	}

	service := options.services.Get(call.Service)
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

// errorHandler constructs a new error object from the given parameter map and codec
func errorHandler(ctx *broker.Context, node *specs.Node, constructor codec.Constructor, handle *specs.OnError, options Options) (*flow.OnError, error) {
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
