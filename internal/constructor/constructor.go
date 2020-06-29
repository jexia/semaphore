package constructor

import (
	"github.com/jexia/maestro/internal/checks"
	"github.com/jexia/maestro/internal/codec"
	"github.com/jexia/maestro/internal/compare"
	"github.com/jexia/maestro/internal/dependencies"
	"github.com/jexia/maestro/internal/flow"
	"github.com/jexia/maestro/internal/functions"
	"github.com/jexia/maestro/internal/references"
	"github.com/jexia/maestro/pkg/core/api"
	"github.com/jexia/maestro/pkg/core/instance"
	"github.com/jexia/maestro/pkg/core/logger"
	"github.com/jexia/maestro/pkg/core/trace"
	"github.com/jexia/maestro/pkg/metadata"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/labels"
	"github.com/jexia/maestro/pkg/specs/template"
	"github.com/jexia/maestro/pkg/specs/types"
	"github.com/jexia/maestro/pkg/transport"
	"github.com/sirupsen/logrus"
)

// Specs construct a specs manifest from the given options
func Specs(ctx instance.Context, mem functions.Collection, options api.Options) (*api.Collection, error) {
	collection, err := CollectSpecs(ctx, options)
	if err != nil {
		return nil, err
	}

	ConstructErrorHandle(collection.Flows)

	err = checks.ManifestDuplicates(ctx, collection.Flows)
	if err != nil {
		return nil, err
	}

	err = references.DefineManifest(ctx, collection.Services, collection.Schema, collection.Flows)
	if err != nil {
		return nil, err
	}

	err = functions.PrepareManifestFunctions(ctx, mem, options.Functions, collection.Flows)
	if err != nil {
		return nil, err
	}

	dependencies.ResolveReferences(ctx, collection.Flows)

	err = compare.ManifestTypes(ctx, collection.Services, collection.Schema, collection.Flows)
	if err != nil {
		return nil, err
	}

	err = dependencies.ResolveManifest(ctx, collection.Flows)
	if err != nil {
		return nil, err
	}

	if options.AfterConstructor != nil {
		err = options.AfterConstructor(ctx, collection)
		if err != nil {
			return nil, err
		}
	}

	return collection, nil
}

// FlowManager constructs the flow managers from the given specs manifest
func FlowManager(ctx instance.Context, mem functions.Collection, services *specs.ServicesManifest, endpoints *specs.EndpointsManifest, flows *specs.FlowsManifest, options api.Options) ([]*transport.Endpoint, error) {
	results := make([]*transport.Endpoint, len(endpoints.Endpoints))

	ctx.Logger(logger.Core).WithField("endpoints", endpoints.Endpoints).Debug("constructing endpoints")

	for index, endpoint := range endpoints.Endpoints {
		manager := flows.GetFlow(endpoint.Flow)
		if manager == nil {
			continue
		}

		nodes := make([]*flow.Node, len(manager.GetNodes()))

		for index, node := range manager.GetNodes() {
			condition := Condition(ctx, mem, node.Condition)

			caller, err := Call(ctx, mem, services, flows, node, node.Call, options, manager)
			if err != nil {
				return nil, err
			}

			rollback, err := Call(ctx, mem, services, flows, node, node.Rollback, options, manager)
			if err != nil {
				return nil, err
			}

			nodes[index] = flow.NewNode(ctx, node, condition, caller, rollback, &flow.NodeMiddleware{
				BeforeDo:       options.BeforeNodeDo,
				AfterDo:        options.AfterNodeDo,
				BeforeRollback: options.BeforeNodeRollback,
				AfterRollback:  options.AfterNodeRollback,
			})
		}

		forward, err := Forward(services, flows, manager.GetForward(), options)
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

	err := Listeners(results, options)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// Condition constructs a new flow condition of the given specs
func Condition(ctx instance.Context, mem functions.Collection, condition *specs.Condition) *flow.Condition {
	if condition == nil {
		return nil
	}

	stack := mem[condition.Params]
	return flow.NewCondition(stack, condition)
}

// Call constructs a flow caller for the given node call.
func Call(ctx instance.Context, mem functions.Collection, services *specs.ServicesManifest, flows *specs.FlowsManifest, node *specs.Node, call *specs.Call, options api.Options, manager specs.FlowResourceManager) (flow.Call, error) {
	if call == nil {
		return nil, nil
	}

	if call.Service != "" {
		return NewServiceCall(ctx, mem, services, flows, node, call, options, manager)
	}

	request, err := Request(ctx, node, mem, nil, call.Request)
	if err != nil {
		return nil, err
	}

	response, err := Request(ctx, node, mem, nil, call.Response)
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

	unexpected, err := Error(ctx, node, mem, codec, node.OnError)
	if err != nil {
		return nil, err
	}

	request, err := Request(ctx, node, mem, codec, call.Request)
	if err != nil {
		return nil, err
	}

	response, err := Request(ctx, node, mem, codec, call.Response)
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

// Request constructs a new request from the given parameter map and codec
func Request(ctx instance.Context, node *specs.Node, mem functions.Collection, constructor codec.Constructor, params *specs.ParameterMap) (*flow.Request, error) {
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

// Error constructs a new error object from the given parameter map and codec
func Error(ctx instance.Context, node *specs.Node, mem functions.Collection, constructor codec.Constructor, err *specs.OnError) (*flow.OnError, error) {
	if err == nil {
		return nil, nil
	}

	var codec codec.Manager
	var meta *metadata.Manager
	var stack functions.Stack

	if err.Response != nil && constructor != nil {
		params := err.Response

		// TODO: check if I would like props to be defined like this
		manager, err := constructor.New(template.JoinPath(node.Name, template.ErrorResource), params)
		if err != nil {
			return nil, err
		}

		codec = manager
		stack = mem[params]
		meta = metadata.NewManager(ctx, node.Name, params.Header)
	}

	return flow.NewOnError(stack, codec, meta, err.Status, err.Message), nil
}

// Forward constructs a flow caller for the given call.
func Forward(services *specs.ServicesManifest, flows *specs.FlowsManifest, call *specs.Call, options api.Options) (*transport.Forward, error) {
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

// Listeners constructs the listeners from the given collection of endpoints
func Listeners(endpoints []*transport.Endpoint, options api.Options) error {
	collections := make(map[string][]*transport.Endpoint, len(options.Listeners))

	options.Ctx.Logger(logger.Core).WithField("endpoints", endpoints).Debug("constructing listeners")

	for _, endpoint := range endpoints {
		if endpoint == nil {
			continue
		}

		options.Ctx.Logger(logger.Core).WithFields(logrus.Fields{
			"flow":     endpoint.Flow.GetName(),
			"listener": endpoint.Listener,
		}).Info("Preparing endpoint")

		listener := options.Listeners.Get(endpoint.Listener)
		if listener == nil {
			options.Ctx.Logger(logger.Core).WithFields(logrus.Fields{
				"listener": endpoint.Listener,
			}).Error("Listener not found")

			return trace.New(trace.WithMessage("unknown listener %s", endpoint.Listener))
		}

		collections[endpoint.Listener] = append(collections[endpoint.Listener], endpoint)
	}

	for key, collection := range collections {
		options.Ctx.Logger(logger.Core).WithField("listener", key).Debug("applying listener handles")

		listener := options.Listeners.Get(key)
		err := listener.Handle(options.Ctx, collection, options.Codec)
		if err != nil {
			options.Ctx.Logger(logger.Core).WithFields(logrus.Fields{
				"listener": listener.Name(),
				"err":      err,
			}).Error("Listener returned an error")

			return err
		}
	}

	return nil
}

// DefaultOnError sets the default values for not defined properties
func DefaultOnError(err *specs.OnError) {
	if err == nil {
		err = &specs.OnError{}
	}

	if err.Status == nil {
		err.Status = &specs.Property{
			Type:  types.Int64,
			Label: labels.Optional,
			Reference: &specs.PropertyReference{
				Resource: "error",
				Path:     "status",
			},
		}
	}

	if err.Message == nil {
		err.Message = &specs.Property{
			Type:  types.String,
			Label: labels.Optional,
			Reference: &specs.PropertyReference{
				Resource: "error",
				Path:     "message",
			},
		}
	}
}

// MergeOnError merges the right on error specs into the left on error
func MergeOnError(left *specs.OnError, right *specs.OnError) {
	if left == nil || right == nil {
		return
	}

	if left.Message == nil {
		left.Message = right.Message.Clone()
	}

	if left.Status == nil {
		left.Status = right.Status.Clone()
	}

	if len(left.Params) == 0 {
		left.Params = make(map[string]*specs.Property, len(right.Params))

		for key, param := range right.Params {
			left.Params[key] = param.Clone()
		}
	}

	if left.Response == nil {
		left.Response = right.Response.Clone()
	}
}

// ConstructErrorHandle clones any previously defined error objects or error handles
func ConstructErrorHandle(manifest *specs.FlowsManifest) {
	for _, flow := range manifest.Flows {
		DefaultOnError(flow.OnError)

		if flow.OnError.Response == nil {
			flow.OnError.Response = manifest.Error.Clone()
		}

		for _, node := range flow.Nodes {
			if node.OnError == nil {
				node.OnError = flow.OnError.Clone()
				continue
			}

			MergeOnError(node.OnError, flow.OnError)
		}
	}

	for _, proxy := range manifest.Proxy {
		DefaultOnError(proxy.OnError)

		if proxy.OnError.Response == nil {
			proxy.OnError.Response = manifest.Error.Clone()
		}

		for _, node := range proxy.Nodes {
			if node.OnError == nil {
				node.OnError = proxy.OnError.Clone()
				continue
			}

			MergeOnError(node.OnError, proxy.OnError)
		}
	}
}
