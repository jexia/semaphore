package conditions

import (
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/logger"
	"github.com/jexia/maestro/pkg/specs"
)

// ResolveExpressions resolves the available expressions within the given flows
func ResolveExpressions(ctx instance.Context, manifest *specs.FlowsManifest) error {
	for _, flow := range manifest.Flows {
		for _, node := range flow.Nodes {
			err := ResolveNodeExpressions(ctx, node)
			if err != nil {
				return err
			}
		}
	}

	for _, proxy := range manifest.Proxy {
		for _, node := range proxy.Nodes {
			err := ResolveNodeExpressions(ctx, node)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ResolveNodeExpressions resolves the expressions available inside the given node
func ResolveNodeExpressions(ctx instance.Context, node *specs.Node) error {
	ctx.Logger(logger.Core).WithField("node", node.Name).Debug("resolving condition expressions")

	if node == nil {
		return nil
	}

	// TODO: resolve node expression

	return nil
}
