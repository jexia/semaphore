package constructor

import (
	"github.com/jexia/maestro/pkg/codec"
	"github.com/jexia/maestro/pkg/conditions"
	"github.com/jexia/maestro/pkg/flow"
	"github.com/jexia/maestro/pkg/functions"
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/logger"
	"github.com/jexia/maestro/pkg/metadata"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/checks"
	"github.com/jexia/maestro/pkg/specs/compare"
	"github.com/jexia/maestro/pkg/specs/dependencies"
	"github.com/jexia/maestro/pkg/specs/references"
	"github.com/jexia/maestro/pkg/specs/trace"
	"github.com/jexia/maestro/pkg/transport"
	"github.com/sirupsen/logrus"
)

// Specs construct a specs manifest from the given options
func Specs(ctx instance.Context, mem functions.Collection, options Options) (*Collection, error) {
	collection, err := CollectSpecs(ctx, options)
	if err != nil {
		return nil, err
	}

	err = references.DefineManifest(ctx, collection.Services, collection.Schema, collection.Flows)
	if err != nil {
		return nil, err
	}

	err = checks.ManifestDuplicates(ctx, collection.Flows)
	if err != nil {
		return nil, err
	}

	err = functions.PrepareManifestFunctions(ctx, mem, options.Functions, collection.Flows)
	if err != nil {
		return nil, err
	}

	err = compare.ManifestTypes(ctx, collection.Services, collection.Schema, collection.Flows)
	if err != nil {
		return nil, err
	}

	dependencies.ResolveReferences(ctx, collection.Flows)

	err = conditions.ResolveExpressions(ctx, collection.Flows)
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
func FlowManager(ctx instance.Context, mem functions.Collection, services *specs.ServicesManifest, endpoints *specs.EndpointsManifest, flows *specs.FlowsManifest, options Options) ([]*transport.Endpoint, error) {
	results := make([]*transport.Endpoint, len(endpoints.Endpoints))

	ctx.Logger(logger.Core).WithField("endpoints", endpoints.Endpoints).Debug("constructing endpoints")

	for index, endpoint := range endpoints.Endpoints {
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

			nodes[index] = flow.NewNode(ctx, node, caller, rollback, &flow.NodeMiddleware{
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

		result.Forward = forward
		result.Flow = flow.NewManager(ctx, manager.GetName(), nodes, &flow.ManagerMiddleware{
			BeforeDo:       options.BeforeManagerDo,
			AfterDo:        options.AfterManagerDo,
			BeforeRollback: options.BeforeManagerRollback,
			AfterRollback:  options.AfterManagerRollback,
		})

		results[index] = result
	}

	err := Listeners(results, options)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// Call constructs a flow caller for the given node call.
func Call(ctx instance.Context, mem functions.Collection, services *specs.ServicesManifest, flows *specs.FlowsManifest, node *specs.Node, call *specs.Call, options Options, manager specs.FlowResourceManager) (flow.Call, error) {
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

	request, err := Request(ctx, node, mem, codec, call.Request)
	if err != nil {
		return nil, err
	}

	response, err := Request(ctx, node, mem, codec, call.Response)
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
