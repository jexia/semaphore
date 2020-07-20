package core

import (
	"github.com/jexia/semaphore/pkg/core/api"
	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/core/logger"
	"github.com/jexia/semaphore/pkg/flow"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/transport"
)

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
			condition := NewCondition(ctx, mem, node.Condition)

			caller, err := NewNodeCall(ctx, mem, services, flows, node, node.Call, options, manager)
			if err != nil {
				return nil, err
			}

			rollback, err := NewNodeCall(ctx, mem, services, flows, node, node.Rollback, options, manager)
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

		forward, err := NewForward(services, flows, manager.GetForward(), options)
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

	err := NewListeners(results, options)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// NewNodeCall constructs a flow caller for the given node call.
func NewNodeCall(ctx instance.Context, mem functions.Collection, services *specs.ServicesManifest, flows *specs.FlowsManifest, node *specs.Node, call *specs.Call, options api.Options, manager specs.FlowResourceManager) (flow.Call, error) {
	if call == nil {
		return nil, nil
	}

	if call.Service != "" {
		return NewServiceCall(ctx, mem, services, flows, node, call, options, manager)
	}

	request, err := NewRequest(ctx, node, mem, nil, call.Request)
	if err != nil {
		return nil, err
	}

	response, err := NewRequest(ctx, node, mem, nil, call.Response)
	if err != nil {
		return nil, err
	}

	caller := flow.NewCall(ctx, node, &flow.CallOptions{
		Request:  request,
		Response: response,
	})

	return caller, nil
}
