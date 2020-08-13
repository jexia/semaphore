package dependencies

import (
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/broker/trace"
	"github.com/jexia/semaphore/pkg/specs"
)

// ResolveFlows resolves all dependencies inside the given manifest
func ResolveFlows(ctx *broker.Context, flows specs.FlowListInterface) error {
	logger.Info(ctx, "resolving flow dependencies")

	for _, flow := range flows {
		for _, node := range flow.GetNodes() {
			err := ResolveNode(flow, node, make(specs.Dependencies))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ResolveNode resolves the given call dependencies and attempts to detect any circular dependencies
func ResolveNode(manager specs.FlowInterface, node *specs.Node, unresolved specs.Dependencies) error {
	if len(node.DependsOn) == 0 {
		return nil
	}

	unresolved[node.ID] = node

	for edge := range node.DependsOn {
		// Remove any self references
		if edge == node.ID {
			delete(unresolved, edge)
			delete(node.DependsOn, edge)
			continue
		}

		_, unresolv := unresolved[edge]
		if unresolv {
			return trace.New(trace.WithMessage("Resource dependencies, circular dependency detected: %s.%s <-> %s.%s", manager.GetName(), node.ID, manager.GetName(), edge))
		}

		result := manager.GetNodes().Get(edge)
		if result == nil {
			continue
		}

		err := ResolveNode(manager, result, unresolved)
		if err != nil {
			return err
		}

		node.DependsOn[edge] = result
	}

	delete(unresolved, node.ID)

	return nil
}
