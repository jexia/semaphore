package functions

import (
	"github.com/jexia/semaphore/v2/pkg/broker"
	"github.com/jexia/semaphore/v2/pkg/references"
	"github.com/jexia/semaphore/v2/pkg/specs"
)

// DefineFunctions defined all properties within the given functions
func DefineFunctions(ctx *broker.Context, functions Stack, node *specs.Node, flow specs.FlowInterface) error {
	if functions == nil {
		return nil
	}

	for _, function := range functions {
		if function.Arguments != nil {
			for _, arg := range function.Arguments {
				references.ResolveProperty(ctx, node, arg, flow)
			}
		}

		if function.Returns == nil {
			continue
		}

		references.ResolveProperty(ctx, node, function.Returns, flow)
	}

	return nil
}
