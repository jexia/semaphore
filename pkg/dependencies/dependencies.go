package dependencies

import (
	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/core/logger"
	"github.com/jexia/semaphore/pkg/core/trace"
	"github.com/jexia/semaphore/pkg/specs"
)

// ResolveManifest resolves all dependencies inside the given manifest
func ResolveManifest(ctx instance.Context, manifest *specs.FlowsManifest) error {
	ctx.Logger(logger.Core).Info("Resolving manifest dependencies")

	for _, flow := range manifest.Flows {
		for _, node := range flow.Nodes {
			err := ResolveNode(flow, node, make(map[string]*specs.Node))
			if err != nil {
				return err
			}
		}
	}

	for _, proxy := range manifest.Proxy {
		for _, node := range proxy.Nodes {
			err := ResolveNode(proxy, node, make(map[string]*specs.Node))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ResolveNode resolves the given call dependencies and attempts to detect any circular dependencies
func ResolveNode(manager specs.FlowResourceManager, node *specs.Node, unresolved map[string]*specs.Node) error {
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
			return trace.New(trace.WithMessage("Circular dependency detected: %s.%s <-> %s.%s", manager.GetName(), node.ID, manager.GetName(), edge))
		}

		result := FindNode(manager, edge)
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

// FindNode attempts to find the given node inside the given flow manager
func FindNode(manager specs.FlowResourceManager, node string) *specs.Node {
	for _, inner := range manager.GetNodes() {
		if inner.ID == node {
			return inner
		}
	}

	return nil
}
