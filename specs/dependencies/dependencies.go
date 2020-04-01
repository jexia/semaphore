package dependencies

import (
	"github.com/jexia/maestro/instance"
	"github.com/jexia/maestro/logger"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/trace"
)

// ResolveManifest resolves all dependencies inside the given manifest
func ResolveManifest(ctx instance.Context, manifest *specs.Manifest) error {
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
func ResolveNode(manager specs.FlowManager, node *specs.Node, unresolved map[string]*specs.Node) error {
	if len(node.DependsOn) == 0 {
		return nil
	}

	unresolved[node.Name] = node

lookup:
	for edge := range node.DependsOn {
		_, unresolv := unresolved[edge]
		if unresolv {
			return trace.New(trace.WithMessage("Circular dependency detected: %s.%s <-> %s.%s", manager.GetName(), node.Name, manager.GetName(), edge))
		}

		for _, inner := range manager.GetNodes() {
			if inner.Name == edge {
				err := ResolveNode(manager, inner, unresolved)
				if err != nil {
					return err
				}

				node.DependsOn[edge] = inner
				continue lookup
			}
		}
	}

	delete(unresolved, node.Name)

	return nil
}
